package response

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

type NotificationResp struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	TaskID    string     `json:"task_id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Read      bool       `json:"read"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

type ListNotificationsResp struct {
	Notifications []NotificationResp `json:"notifications"`
	Total         int                `json:"total"`
}

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

	return ListNotificationsResp{Notifications: items, Total: total}
}
