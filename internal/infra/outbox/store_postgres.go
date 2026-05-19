package outbox

import (
	"context"
	"fmt"
	"time"

	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// Store persists and claims outbox events.
type Store struct {
	executor postgres.Executor
}

// NewStore creates a Postgres outbox store.
func NewStore(executor postgres.Executor) *Store {
	return &Store{executor: executor}
}

// Add stores an event. Call this with a transaction executor when the event must commit with business data.
func (s *Store) Add(ctx context.Context, event *appports.OutboxEvent) error {
	record := Event{
		ID:            event.ID,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID,
		EventType:     event.EventType,
		Payload:       event.Payload,
		Headers:       event.Headers,
		AvailableAt:   event.AvailableAt,
	}

	return s.addRecord(ctx, &record)
}

func (s *Store) addRecord(ctx context.Context, event *Event) error {
	now := time.Now().UTC()

	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	if event.AvailableAt.IsZero() {
		event.AvailableAt = now
	}

	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}

	headers, err := json.Marshal(event.Headers)
	if err != nil {
		return fmt.Errorf("outbox Store - Add - marshal headers: %w", err)
	}

	_, err = s.executor.Exec(ctx, `
		INSERT INTO outbox_events (
			id, aggregate_type, aggregate_id, event_type, payload, headers, status, attempts, available_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		event.ID,
		event.AggregateType,
		event.AggregateID,
		event.EventType,
		event.Payload,
		headers,
		StatusPending,
		event.Attempts,
		event.AvailableAt,
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("outbox Store - Add - Exec: %w", err)
	}

	return nil
}

// ClaimPending atomically claims pending events and stale publishing locks for one relay worker.
func (s *Store) ClaimPending(ctx context.Context, limit int, lockTimeout time.Duration) ([]Event, error) {
	if limit <= 0 {
		limit = 20
	}

	if lockTimeout <= 0 {
		lockTimeout = defaultLockTimeout
	}

	rows, err := s.executor.Query(ctx, `
		UPDATE outbox_events
		SET status = $1, locked_at = now(), attempts = attempts + 1
		WHERE id IN (
			SELECT id
			FROM outbox_events
			WHERE (status = $2 AND available_at <= now())
			   OR (status = $3 AND locked_at IS NOT NULL AND locked_at <= now() - make_interval(secs => $4))
			ORDER BY created_at
			LIMIT $5
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, aggregate_type, aggregate_id, event_type, payload, headers, status, attempts,
		          available_at, created_at, locked_at, published_at, last_error
	`, StatusPublishing, StatusPending, StatusPublishing, int(lockTimeout.Seconds()), limit)
	if err != nil {
		return nil, fmt.Errorf("outbox Store - ClaimPending - Query: %w", err)
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("outbox Store - ClaimPending - rows.Err: %w", err)
	}

	return events, nil
}

// MarkPublished marks an event as published.
func (s *Store) MarkPublished(ctx context.Context, id string) error {
	_, err := s.executor.Exec(ctx, `
		UPDATE outbox_events
		SET status = $1, published_at = now(), last_error = NULL
		WHERE id = $2
	`, StatusPublished, id)
	if err != nil {
		return fmt.Errorf("outbox Store - MarkPublished - Exec: %w", err)
	}

	return nil
}

// MarkFailed schedules an event retry or marks it failed after max attempts.
func (s *Store) MarkFailed(ctx context.Context, id string, publishErr error, maxAttempts int) error {
	if maxAttempts <= 0 {
		maxAttempts = 10
	}

	_, err := s.executor.Exec(ctx, `
		UPDATE outbox_events
		SET status = CASE WHEN attempts >= $1 THEN $2 ELSE $3 END,
		    available_at = now() + make_interval(secs => LEAST(60, attempts * attempts)),
		    last_error = $4
		WHERE id = $5
	`, maxAttempts, StatusFailed, StatusPending, publishErr.Error(), id)
	if err != nil {
		return fmt.Errorf("outbox Store - MarkFailed - Exec: %w", err)
	}

	return nil
}
