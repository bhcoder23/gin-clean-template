// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errNilPool = errors.New("postgres pool is nil")

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	logger       logger.Interface

	Pool *pgxpool.Pool
}

// Executor is the minimal query contract shared by pgxpool.Pool and pgx.Tx.
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = safeIntToInt32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		pg.info("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) info(message string, args ...any) {
	if p.logger == nil {
		return
	}

	p.logger.Info(message, args...)
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// Ping checks database readiness.
func (p *Postgres) Ping(ctx context.Context) error {
	if p.Pool == nil {
		return errNilPool
	}

	return p.Pool.Ping(ctx)
}

// WithinTx runs fn in a transaction and passes the transaction executor to it.
func (p *Postgres) WithinTx(ctx context.Context, fn func(context.Context, Executor) error) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres - WithinTx - Begin: %w", err)
	}

	if err = fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("postgres - WithinTx - rollback after %w: %w", err, rollbackErr)
		}

		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("postgres - WithinTx - Commit: %w", err)
	}

	return nil
}

func safeIntToInt32(v int) int32 {
	if v <= 0 {
		return 1
	}

	if v > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(v)
}
