package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
)

// NotificationUsecase coordinates notification application workflows.
type NotificationUsecase struct {
	repo appports.NotificationRepo
}

// New -.
func New(r appports.NotificationRepo) *NotificationUsecase {
	return &NotificationUsecase{repo: r}
}

// List -.
func (uc *NotificationUsecase) List(ctx context.Context, userID string, unreadOnly *bool, limit, offset int) ([]domain.Notification, int, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	notifications, total, err := uc.repo.List(ctx, userID, appports.NotificationFilter{
		UnreadOnly: unreadOnly,
		Limit:      uint64(limit),
		Offset:     uint64(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationUsecase - List - uc.repo.List: %w", err)
	}

	return notifications, total, nil
}

// MarkRead -.
func (uc *NotificationUsecase) MarkRead(ctx context.Context, userID, notificationID string) (domain.Notification, error) {
	notification, err := uc.repo.GetByID(ctx, userID, notificationID)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("NotificationUsecase - MarkRead - uc.repo.GetByID: %w", err)
	}

	if notification.Read {
		return notification, nil
	}

	now := time.Now().UTC()
	notification.Read = true
	notification.ReadAt = &now

	if err = uc.repo.Update(ctx, &notification); err != nil {
		return domain.Notification{}, fmt.Errorf("NotificationUsecase - MarkRead - uc.repo.Update: %w", err)
	}

	return notification, nil
}
