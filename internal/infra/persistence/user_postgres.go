package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

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

type userRow struct {
	id           string
	username     string
	email        string
	passwordHash string
	createdAt    time.Time
	updatedAt    time.Time
}

func newUserRow(user domain.User) userRow {
	return userRow{
		id:           user.ID,
		username:     user.Username,
		email:        user.Email,
		passwordHash: user.PasswordHash,
		createdAt:    user.CreatedAt,
		updatedAt:    user.UpdatedAt,
	}
}

func (r userRow) toDomain() domain.User {
	return domain.User{
		ID:           r.id,
		Username:     r.username,
		Email:        r.email,
		PasswordHash: r.passwordHash,
		CreatedAt:    r.createdAt,
		UpdatedAt:    r.updatedAt,
	}
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
	row := newUserRow(*user)

	sql, args, err := r.builder.
		Insert("users").
		Columns("id, username, email, password_hash, created_at, updated_at").
		Values(row.id, row.username, row.email, row.passwordHash, row.createdAt, row.updatedAt).
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

	var row userRow

	err = r.executor.QueryRow(ctx, sql, args...).
		Scan(&row.id, &row.username, &row.email, &row.passwordHash, &row.createdAt, &row.updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("UserRepo - getUser - r.Pool.QueryRow: %w", err)
	}

	return row.toDomain(), nil
}
