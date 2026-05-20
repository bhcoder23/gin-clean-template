package response

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

// NotificationResp documents the REST notification response shape.
type NotificationResp struct {
	ID        string     `example:"550e8400-e29b-41d4-a716-446655440000" json:"id"`
	UserID    string     `example:"550e8400-e29b-41d4-a716-446655440000" json:"user_id"`
	TaskID    string     `example:"550e8400-e29b-41d4-a716-446655440000" json:"task_id"`
	Type      string     `example:"task_created"                        json:"type"`
	Title     string     `example:"Task created"                        json:"title"`
	Body      string     `example:"Task was created."                    json:"body"`
	Read      bool       `example:"false"                               json:"read"`
	CreatedAt time.Time  `example:"2026-01-01T00:00:00Z"                json:"created_at"`
	ReadAt    *time.Time `example:"2026-01-01T00:00:00Z"                json:"read_at,omitempty"`
} // @name v1.NotificationResp

// ListNotificationsResp -.
type ListNotificationsResp struct {
	Notifications []NotificationResp `json:"notifications"`
	Total         int                `example:"42" json:"total"`
} // @name v1.ListNotificationsResp

func NewNotificationResp(notification domain.Notification) NotificationResp {
	return NotificationResp{
		ID:        notification.ID,
		UserID:    notification.UserID,
		TaskID:    notification.TaskID,
		Type:      string(notification.Type),
		Title:     notification.Title,
		Body:      notification.Body,
		Read:      notification.Read,
		CreatedAt: notification.CreatedAt,
		ReadAt:    notification.ReadAt,
	}
}

func NewListNotificationsResp(notifications []domain.Notification, total int) ListNotificationsResp {
	items := make([]NotificationResp, 0, len(notifications))
	for _, notification := range notifications {
		items = append(items, NewNotificationResp(notification))
	}

	return ListNotificationsResp{
		Notifications: items,
		Total:         total,
	}
}
