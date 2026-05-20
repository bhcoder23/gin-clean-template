package domain

import (
	"errors"
	"slices"
	"time"
)

// TaskStatus -.
type TaskStatus string

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
	ID          string
	UserID      string
	Title       string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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
