package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTaskNotFound         = errors.New("task not found")
	ErrTaskTitleRequired    = errors.New("task title is required")
	ErrTaskCompleted        = errors.New("completed task cannot be modified")
	ErrInvalidTransition    = errors.New("invalid status transition")
	ErrNotificationNotFound = errors.New("notification not found")
)
