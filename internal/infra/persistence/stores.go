package persistence

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/infra/outbox"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
)

// Repositories exposes repositories bound to the same executor.
type Repositories struct {
	executor postgres.Executor
}

// NewRepositories creates repositories backed by the normal database pool.
func NewRepositories(pg *postgres.Postgres) Repositories {
	return NewRepositoriesWithExecutor(pg.Pool)
}

// NewRepositoriesWithExecutor creates repositories backed by a pool or transaction executor.
func NewRepositoriesWithExecutor(executor postgres.Executor) Repositories {
	return Repositories{
		executor: executor,
	}
}

// Users returns a user repository.
func (r Repositories) Users() appports.UserRepo {
	return NewUserRepoWithExecutor(r.executor)
}

// Tasks returns a task repository.
func (r Repositories) Tasks() appports.TaskRepo {
	return NewTaskRepoWithExecutor(r.executor)
}

// Notifications returns a notification repository.
func (r Repositories) Notifications() appports.NotificationRepo {
	return NewNotificationRepoWithExecutor(r.executor)
}

// Outbox returns an outbox repository.
func (r Repositories) Outbox() appports.OutboxStore {
	return outbox.NewStore(r.executor)
}

// Transactor creates transaction-scoped repositories.
type Transactor struct {
	pg *postgres.Postgres
}

// NewTransactor returns a transaction helper for use cases that need atomic multi-repo work.
func NewTransactor(pg *postgres.Postgres) *Transactor {
	return &Transactor{pg: pg}
}

// WithinTx runs fn with repositories bound to one transaction.
func (t *Transactor) WithinTx(ctx context.Context, fn func(context.Context, appports.RepoProvider) error) error {
	return t.pg.WithinTx(ctx, func(txCtx context.Context, executor postgres.Executor) error {
		return fn(txCtx, NewRepositoriesWithExecutor(executor))
	})
}
