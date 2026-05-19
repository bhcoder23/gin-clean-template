package persistence

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/bhcoder23/gin-clean-template/internal/infra/outbox"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
)

// Stores exposes repositories bound to the same executor.
type Stores struct {
	builder  sq.StatementBuilderType
	executor postgres.Executor
}

// NewStores creates repositories backed by the normal database pool.
func NewStores(pg *postgres.Postgres) Stores {
	return NewStoresWithExecutor(pg.Builder, pg.Pool)
}

// NewStoresWithExecutor creates repositories backed by a pool or transaction executor.
func NewStoresWithExecutor(builder sq.StatementBuilderType, executor postgres.Executor) Stores {
	return Stores{
		builder:  builder,
		executor: executor,
	}
}

// Users returns a user repository.
func (s Stores) Users() appports.UserStore {
	return NewUserRepoWithExecutor(s.builder, s.executor)
}

// Tasks returns a task repository.
func (s Stores) Tasks() appports.TaskStore {
	return NewTaskRepoWithExecutor(s.builder, s.executor)
}

// Notifications returns a notification repository.
func (s Stores) Notifications() appports.NotificationStore {
	return NewNotificationRepoWithExecutor(s.builder, s.executor)
}

// Outbox returns an outbox repository.
func (s Stores) Outbox() appports.OutboxStore {
	return outbox.NewStore(s.executor)
}

// Transactor creates transaction-scoped stores.
type Transactor struct {
	pg *postgres.Postgres
}

// NewTransactor returns a transaction helper for use cases that need atomic multi-repo work.
func NewTransactor(pg *postgres.Postgres) *Transactor {
	return &Transactor{pg: pg}
}

// WithinTx runs fn with repositories bound to one transaction.
func (t *Transactor) WithinTx(ctx context.Context, fn func(context.Context, appports.StoreProvider) error) error {
	return t.pg.WithinTx(ctx, func(txCtx context.Context, executor postgres.Executor) error {
		return fn(txCtx, NewStoresWithExecutor(t.pg.Builder, executor))
	})
}
