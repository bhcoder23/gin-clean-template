package response

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

// TaskResp documents the REST task response shape.
type TaskResp struct {
	ID          string    `example:"550e8400-e29b-41d4-a716-446655440000" json:"id"`
	UserID      string    `example:"550e8400-e29b-41d4-a716-446655440000" json:"user_id"`
	Title       string    `example:"My task"                              json:"title"`
	Description string    `example:"Task description"                     json:"description"`
	Status      string    `example:"todo"                                 json:"status"`
	CreatedAt   time.Time `example:"2026-01-01T00:00:00Z"                 json:"created_at"`
	UpdatedAt   time.Time `example:"2026-01-01T00:00:00Z"                 json:"updated_at"`
} // @name v1.TaskResp

// ListTasksResp -.
type ListTasksResp struct {
	Tasks []TaskResp `json:"tasks"`
	Total int        `example:"42" json:"total"`
} // @name v1.ListTasksResp

func NewTaskResp(task *domain.Task) TaskResp {
	return TaskResp{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func NewListTasksResp(tasks []domain.Task, total int) ListTasksResp {
	items := make([]TaskResp, 0, len(tasks))
	for i := range tasks {
		items = append(items, NewTaskResp(&tasks[i]))
	}

	return ListTasksResp{
		Tasks: items,
		Total: total,
	}
}
