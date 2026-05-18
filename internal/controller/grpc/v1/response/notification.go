package response

import (
	"math"

	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/entity"
)

// NewNotificationResponse -.
func NewNotificationResponse(notification *entity.Notification) *v1.NotificationResponse {
	readAt := ""
	if notification.ReadAt != nil {
		readAt = notification.ReadAt.Format("2006-01-02T15:04:05Z")
	}

	return &v1.NotificationResponse{
		Id:        notification.ID,
		UserId:    notification.UserID,
		TaskId:    notification.TaskID,
		Type:      string(notification.Type),
		Title:     notification.Title,
		Body:      notification.Body,
		Read:      notification.Read,
		CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z"),
		ReadAt:    readAt,
	}
}

// NewListNotificationsResponse -.
func NewListNotificationsResponse(notifications []entity.Notification, total int) *v1.ListNotificationsResponse {
	pbNotifications := make([]*v1.NotificationResponse, len(notifications))
	for i := range notifications {
		pbNotifications[i] = NewNotificationResponse(&notifications[i])
	}

	if total > math.MaxInt32 {
		total = math.MaxInt32
	}

	return &v1.ListNotificationsResponse{
		Notifications: pbNotifications,
		Total:         int32(total),
	}
}
