package domain

import "time"

// NotificationType -.
type NotificationType string // @name domain.NotificationType

const (
	NotificationTypeTaskCreated       NotificationType = "task_created"
	NotificationTypeTaskStatusChanged NotificationType = "task_status_changed"
)

// Notification -.
type Notification struct {
	ID        string           `example:"550e8400-e29b-41d4-a716-446655440000" json:"id"`
	UserID    string           `example:"550e8400-e29b-41d4-a716-446655440000" json:"user_id"`
	TaskID    string           `example:"550e8400-e29b-41d4-a716-446655440000" json:"task_id"`
	Type      NotificationType `example:"task_created"                        json:"type"`
	Title     string           `example:"Task created"                        json:"title"`
	Body      string           `example:"Task \"Ship the scaffold\" was created." json:"body"`
	Read      bool             `example:"false"                               json:"read"`
	CreatedAt time.Time        `example:"2026-01-01T00:00:00Z"                json:"created_at"`
	ReadAt    *time.Time       `example:"2026-01-01T00:00:00Z"                json:"read_at,omitempty"`
} // @name domain.Notification

// NotificationList -.
type NotificationList struct {
	Notifications []Notification `json:"notifications"`
	Total         int            `json:"total"`
} // @name domain.NotificationList
