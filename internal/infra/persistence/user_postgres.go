package persistence

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepo -.
type UserRepo struct {
	builder  sq.StatementBuilderType
	executor postgres.Executor
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return NewUserRepoWithExecutor(pg.Builder, pg.Pool)
}

// NewUserRepoWithExecutor creates a repository bound to a pool or transaction executor.
func NewUserRepoWithExecutor(builder sq.StatementBuilderType, executor postgres.Executor) *UserRepo {
	return &UserRepo{
		builder:  builder,
		executor: executor,
	}
}

// Store -.
func (r *UserRepo) Store(ctx context.Context, user *domain.User) error {
	sql, args, err := r.builder.
		Insert("users").
		Columns("id, username, email, password_hash, created_at, updated_at").
		Values(user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Builder: %w", err)
	}

	_, err = r.executor.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}

		return fmt.Errorf("UserRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetByID -.
func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	return r.getUser(ctx, "id", id)
}

// GetByEmail -.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.getUser(ctx, "email", email)
}

func (r *UserRepo) getUser(ctx context.Context, column, value string) (domain.User, error) {
	sql, args, err := r.builder.
		Select("id, username, email, password_hash, created_at, updated_at").
		From("users").
		Where(sq.Eq{column: value}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("UserRepo - getUser - r.Builder: %w", err)
	}

	var user domain.User

	err = r.executor.QueryRow(ctx, sql, args...).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("UserRepo - getUser - r.Pool.QueryRow: %w", err)
	}

	return user, nil
}
