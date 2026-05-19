package outbox

import (
	"context"
	"errors"
	"time"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"golang.org/x/sync/errgroup"
)

const (
	defaultLockTimeout    = 5 * time.Minute
	defaultPublishTimeout = 5 * time.Second
)

// Publisher publishes claimed outbox events to an external broker.
type Publisher interface {
	Publish(ctx context.Context, event *Event) error
}

// Relay claims and publishes outbox events.
type Relay struct {
	store          *Store
	publisher      Publisher
	logger         logger.Interface
	pollInterval   time.Duration
	batchSize      int
	maxAttempts    int
	lockTimeout    time.Duration
	publishTimeout time.Duration
	stop           chan struct{}
	eg             *errgroup.Group
}

// RelayConfig configures an outbox relay.
type RelayConfig struct {
	PollInterval   time.Duration
	BatchSize      int
	MaxAttempts    int
	LockTimeout    time.Duration
	PublishTimeout time.Duration
}

// NewRelay creates a relay.
func NewRelay(store *Store, publisher Publisher, l logger.Interface, cfg RelayConfig) *Relay {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = time.Second
	}

	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 20
	}

	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 10
	}

	if cfg.LockTimeout <= 0 {
		cfg.LockTimeout = defaultLockTimeout
	}

	if cfg.PublishTimeout <= 0 {
		cfg.PublishTimeout = defaultPublishTimeout
	}

	return &Relay{
		store:          store,
		publisher:      publisher,
		logger:         l,
		pollInterval:   cfg.PollInterval,
		batchSize:      cfg.BatchSize,
		maxAttempts:    cfg.MaxAttempts,
		lockTimeout:    cfg.LockTimeout,
		publishTimeout: cfg.PublishTimeout,
		stop:           make(chan struct{}),
		eg:             new(errgroup.Group),
	}
}

// Start starts the relay loop.
func (r *Relay) Start(ctx context.Context) {
	r.eg.Go(func() error {
		ticker := time.NewTicker(r.pollInterval)
		defer ticker.Stop()

		for {
			if err := r.PublishOnce(ctx); err != nil {
				r.logger.Error(err, "outbox relay - publish once")
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-r.stop:
				return nil
			case <-ticker.C:
			}
		}
	})
}

// Shutdown stops the relay loop.
func (r *Relay) Shutdown() error {
	close(r.stop)

	err := r.eg.Wait()
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}

// PublishOnce claims and publishes one batch.
func (r *Relay) PublishOnce(ctx context.Context) error {
	events, err := r.store.ClaimPending(ctx, r.batchSize, r.lockTimeout)
	if err != nil {
		return err
	}

	for i := range events {
		event := &events[i]

		if err = r.publishEvent(ctx, event); err != nil {
			r.logger.Warn(err, "outbox relay - publish event %s", event.ID)

			if markErr := r.store.MarkFailed(ctx, event.ID, err, r.maxAttempts); markErr != nil {
				return markErr
			}

			continue
		}

		if err := r.store.MarkPublished(ctx, event.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *Relay) publishEvent(ctx context.Context, event *Event) error {
	publishCtx, cancel := context.WithTimeout(ctx, r.publishTimeout)
	defer cancel()

	return r.publisher.Publish(publishCtx, event)
}
