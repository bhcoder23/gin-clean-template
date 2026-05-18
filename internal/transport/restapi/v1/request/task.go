package request

import "github.com/bhcoder23/gin-clean-template/internal/domain"

// CreateTask -.
type CreateTask struct {
	Title       string `example:"My task"          json:"title"       validate:"required,max=255"`
	Description string `example:"Task description" json:"description" validate:"max=1000"`
} // @name v1.CreateTask

// UpdateTask -.
type UpdateTask struct {
	Title       string `example:"Updated task"        json:"title"       validate:"required,max=255"`
	Description string `example:"Updated description" json:"description" validate:"max=1000"`
} // @name v1.UpdateTask

// TransitionTask -.
type TransitionTask struct {
	Status domain.TaskStatus `example:"in_progress" json:"status" validate:"required,oneof=todo in_progress done"`
} // @name v1.TransitionTask
