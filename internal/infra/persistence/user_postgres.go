package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/infra/persistence/sqlc"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepo -.
type UserRepo struct {
	queries *sqlc.Queries
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return NewUserRepoWithExecutor(pg.Pool)
}

// NewUserRepoWithExecutor creates a repository bound to a pool or transaction executor.
func NewUserRepoWithExecutor(executor postgres.Executor) *UserRepo {
	return &UserRepo{
		queries: sqlc.New(executor),
	}
}

// Store -.
func (r *UserRepo) Store(ctx context.Context, user *domain.User) error {
	err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}

		return fmt.Errorf("UserRepo - Store - CreateUser: %w", err)
	}

	return nil
}

// GetByID -.
func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, mapUserReadError("GetByID", err)
	}

	return userToDomain(&row), nil
}

// GetByEmail -.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, mapUserReadError("GetByEmail", err)
	}

	return userToDomain(&row), nil
}

func mapUserReadError(operation string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}

	return fmt.Errorf("UserRepo - %s: %w", operation, err)
}

func userToDomain(row *sqlc.User) domain.User {
	return domain.User{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
