package response

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

type TaskResp struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListTasksResp struct {
	Tasks []TaskResp `json:"tasks"`
	Total int        `json:"total"`
}

func NewTaskResp(task domain.Task) TaskResp {
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
	for _, task := range tasks {
		items = append(items, NewTaskResp(task))
	}

	return ListTasksResp{Tasks: items, Total: total}
}

// DeleteStatus -.
type DeleteStatus struct {
	Status string `json:"status"`
}
