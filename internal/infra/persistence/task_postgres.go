package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/infra/persistence/sqlc"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// TaskRepo -.
type TaskRepo struct {
	queries *sqlc.Queries
}

// NewTaskRepo -.
func NewTaskRepo(pg *postgres.Postgres) *TaskRepo {
	return NewTaskRepoWithExecutor(pg.Pool)
}

// NewTaskRepoWithExecutor creates a repository bound to a pool or transaction executor.
func NewTaskRepoWithExecutor(executor postgres.Executor) *TaskRepo {
	return &TaskRepo{
		queries: sqlc.New(executor),
	}
}

// Store -.
func (r *TaskRepo) Store(ctx context.Context, task *domain.Task) error {
	err := r.queries.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: textParam(task.Description),
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("TaskRepo - Store - CreateTask: %w", err)
	}

	return nil
}

// GetByID -.
func (r *TaskRepo) GetByID(ctx context.Context, userID, taskID string) (domain.Task, error) {
	row, err := r.queries.GetTaskByID(ctx, sqlc.GetTaskByIDParams{
		ID:     taskID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}

		return domain.Task{}, fmt.Errorf("TaskRepo - GetByID - GetTaskByID: %w", err)
	}

	return taskToDomain(&row), nil
}

// List -.
func (r *TaskRepo) List(ctx context.Context, userID string, filter appports.TaskFilter) ([]domain.Task, int, error) {
	params, err := taskListParams(userID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - taskListParams: %w", err)
	}

	total, err := r.queries.CountTasks(ctx, sqlc.CountTasksParams{
		UserID: params.UserID,
		Status: params.Status,
		Query:  params.Query,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - CountTasks: %w", err)
	}

	rows, err := r.queries.ListTasks(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - ListTasks: %w", err)
	}

	return tasksToDomain(rows), int(total), nil
}

// Update -.
func (r *TaskRepo) Update(ctx context.Context, task *domain.Task) error {
	rowsAffected, err := r.queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		Title:       task.Title,
		Description: textParam(task.Description),
		Status:      string(task.Status),
		UpdatedAt:   task.UpdatedAt,
		ID:          task.ID,
		UserID:      task.UserID,
	})
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - UpdateTask: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

// Delete -.
func (r *TaskRepo) Delete(ctx context.Context, userID, taskID string) error {
	rowsAffected, err := r.queries.DeleteTask(ctx, sqlc.DeleteTaskParams{
		ID:     taskID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - DeleteTask: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func taskListParams(userID string, filter appports.TaskFilter) (sqlc.ListTasksParams, error) {
	status := ""
	if filter.Status != nil {
		status = string(*filter.Status)
	}

	limitCount, offsetCount, err := paginationCounts(filter.Limit, filter.Offset)
	if err != nil {
		return sqlc.ListTasksParams{}, err
	}

	return sqlc.ListTasksParams{
		UserID:      userID,
		Status:      status,
		Query:       filter.Query,
		OffsetCount: offsetCount,
		LimitCount:  limitCount,
	}, nil
}

func taskToDomain(row *sqlc.Task) domain.Task {
	return domain.Task{
		ID:          row.ID,
		UserID:      row.UserID,
		Title:       row.Title,
		Description: textValue(row.Description),
		Status:      domain.TaskStatus(row.Status),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func tasksToDomain(rows []sqlc.Task) []domain.Task {
	tasks := make([]domain.Task, 0, len(rows))
	for i := range rows {
		tasks = append(tasks, taskToDomain(&rows[i]))
	}

	return tasks
}
