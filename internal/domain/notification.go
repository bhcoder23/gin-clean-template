package domain

import (
	"errors"
	"time"
)

var ErrNotificationNotFound = errors.New("notification not found")

// NotificationType -.
type NotificationType string

const (
	NotificationTypeTaskCreated       NotificationType = "task_created"
	NotificationTypeTaskStatusChanged NotificationType = "task_status_changed"
)

// Notification -.
type Notification struct {
	ID        string
	UserID    string
	TaskID    string
	Type      NotificationType
	Title     string
	Body      string
	Read      bool
	CreatedAt time.Time
	ReadAt    *time.Time
}
