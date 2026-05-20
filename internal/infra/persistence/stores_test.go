package persistence

import (
	"context"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

type fakeExecutor struct{}

func (fakeExecutor) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (fakeExecutor) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return &fakeRows{}, nil
}

func (fakeExecutor) QueryRow(context.Context, string, ...any) pgx.Row {
	return fakeRow{}
}

type fakeRow struct{}

func (fakeRow) Scan(...any) error {
	return pgx.ErrNoRows
}

var _ postgres.Executor = fakeExecutor{}

func TestRepositoriesCreateRepositoriesBoundToExecutor(t *testing.T) {
	t.Parallel()

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	repos := NewRepositoriesWithExecutor(builder, fakeExecutor{})

	require.NotNil(t, repos.Users())
	require.NotNil(t, repos.Tasks())
	require.NotNil(t, repos.Notifications())
	require.NotNil(t, repos.Outbox())
}
