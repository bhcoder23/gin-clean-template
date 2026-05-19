package domain

import (
	"errors"
	"slices"
	"time"
)

// TaskStatus -.
type TaskStatus string // @name domain.TaskStatus

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskTitleRequired = errors.New("task title is required")
	ErrTaskCompleted     = errors.New("completed task cannot be modified")
	ErrInvalidTransition = errors.New("invalid status transition")
)

// Task -.
type Task struct {
	ID          string     `example:"550e8400-e29b-41d4-a716-446655440000" json:"id"`
	UserID      string     `example:"550e8400-e29b-41d4-a716-446655440000" json:"user_id"`
	Title       string     `example:"My task"                              json:"title"`
	Description string     `example:"Task description"                     json:"description"`
	Status      TaskStatus `example:"todo"                                 json:"status"`
	CreatedAt   time.Time  `example:"2026-01-01T00:00:00Z"                 json:"created_at"`
	UpdatedAt   time.Time  `example:"2026-01-01T00:00:00Z"                 json:"updated_at"`
} // @name domain.Task

// Valid reports whether s is a known task status.
func (s TaskStatus) Valid() bool {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// Transition validates and applies a status transition.
func (t *Task) Transition(newStatus TaskStatus) error {
	validTransitions := map[TaskStatus][]TaskStatus{
		TaskStatusTodo:       {TaskStatusInProgress},
		TaskStatusInProgress: {TaskStatusDone, TaskStatusTodo},
		TaskStatusDone:       {},
	}

	allowed, ok := validTransitions[t.Status]
	if !ok {
		return ErrInvalidTransition
	}

	if slices.Contains(allowed, newStatus) {
		t.Status = newStatus

		return nil
	}

	return ErrInvalidTransition
}
