package persistence

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// TaskRepo -.
type TaskRepo struct {
	builder  sq.StatementBuilderType
	executor postgres.Executor
}

func collectTaskRows(rows pgx.Rows, limit uint64) ([]domain.Task, error) {
	tasks := make([]domain.Task, 0, limit)

	for rows.Next() {
		var t domain.Task

		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("collectTaskRows - rows.Scan: %w", err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("collectTaskRows - rows.Err: %w", err)
	}

	return tasks, nil
}

// NewTaskRepo -.
func NewTaskRepo(pg *postgres.Postgres) *TaskRepo {
	return NewTaskRepoWithExecutor(pg.Builder, pg.Pool)
}

// NewTaskRepoWithExecutor creates a repository bound to a pool or transaction executor.
func NewTaskRepoWithExecutor(builder sq.StatementBuilderType, executor postgres.Executor) *TaskRepo {
	return &TaskRepo{
		builder:  builder,
		executor: executor,
	}
}

// Store -.
func (r *TaskRepo) Store(ctx context.Context, task *domain.Task) error {
	sql, args, err := r.builder.
		Insert("tasks").
		Columns("id, user_id, title, description, status, created_at, updated_at").
		Values(task.ID, task.UserID, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("TaskRepo - Store - r.Builder: %w", err)
	}

	_, err = r.executor.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetByID -.
func (r *TaskRepo) GetByID(ctx context.Context, userID, taskID string) (domain.Task, error) {
	sql, args, err := r.builder.
		Select("id, user_id, title, description, status, created_at, updated_at").
		From("tasks").
		Where(sq.Eq{"id": taskID, "user_id": userID}).
		ToSql()
	if err != nil {
		return domain.Task{}, fmt.Errorf("TaskRepo - GetByID - r.Builder: %w", err)
	}

	var task domain.Task

	err = r.executor.QueryRow(ctx, sql, args...).
		Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}

		return domain.Task{}, fmt.Errorf("TaskRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return task, nil
}

// List -.
func (r *TaskRepo) List(ctx context.Context, userID string, filter appports.TaskFilter) ([]domain.Task, int, error) {
	countBuilder := r.builder.
		Select("COUNT(*)").
		From("tasks").
		Where(sq.Eq{"user_id": userID})

	if filter.Status != nil {
		countBuilder = countBuilder.Where(sq.Eq{"status": *filter.Status})
	}

	if filter.Query != "" {
		countBuilder = countBuilder.Where("title ILIKE ?", "%"+filter.Query+"%")
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - countBuilder: %w", err)
	}

	var total int

	err = r.executor.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - count query: %w", err)
	}

	dataBuilder := r.builder.
		Select("id, user_id, title, description, status, created_at, updated_at").
		From("tasks").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset)

	if filter.Status != nil {
		dataBuilder = dataBuilder.Where(sq.Eq{"status": *filter.Status})
	}

	if filter.Query != "" {
		dataBuilder = dataBuilder.Where("title ILIKE ?", "%"+filter.Query+"%")
	}

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - dataBuilder: %w", err)
	}

	rows, err := r.executor.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	tasks, err := collectTaskRows(rows, filter.Limit)
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - collectTaskRows: %w", err)
	}

	return tasks, total, nil
}

// Update -.
func (r *TaskRepo) Update(ctx context.Context, task *domain.Task) error {
	sql, args, err := r.builder.
		Update("tasks").
		Set("title", task.Title).
		Set("description", task.Description).
		Set("status", task.Status).
		Set("updated_at", task.UpdatedAt).
		Where(sq.Eq{"id": task.ID, "user_id": task.UserID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - r.Builder: %w", err)
	}

	result, err := r.executor.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

// Delete -.
func (r *TaskRepo) Delete(ctx context.Context, userID, taskID string) error {
	sql, args, err := r.builder.
		Delete("tasks").
		Where(sq.Eq{"id": taskID, "user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - r.Builder: %w", err)
	}

	result, err := r.executor.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}
