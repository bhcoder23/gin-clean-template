package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/entity"
	"github.com/bhcoder23/gin-clean-template/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.NotificationRepo
}

// New -.
func New(r repo.NotificationRepo) *UseCase {
	return &UseCase{repo: r}
}

// List -.
func (uc *UseCase) List(ctx context.Context, userID string, unreadOnly *bool, limit, offset int) ([]entity.Notification, int, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	notifications, total, err := uc.repo.List(ctx, userID, repo.NotificationFilter{
		UnreadOnly: unreadOnly,
		Limit:      uint64(limit),
		Offset:     uint64(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationUseCase - List - uc.repo.List: %w", err)
	}

	return notifications, total, nil
}

// MarkRead -.
func (uc *UseCase) MarkRead(ctx context.Context, userID, notificationID string) (entity.Notification, error) {
	notification, err := uc.repo.GetByID(ctx, userID, notificationID)
	if err != nil {
		return entity.Notification{}, fmt.Errorf("NotificationUseCase - MarkRead - uc.repo.GetByID: %w", err)
	}

	if notification.Read {
		return notification, nil
	}

	now := time.Now().UTC()
	notification.Read = true
	notification.ReadAt = &now

	if err = uc.repo.Update(ctx, &notification); err != nil {
		return entity.Notification{}, fmt.Errorf("NotificationUseCase - MarkRead - uc.repo.Update: %w", err)
	}

	return notification, nil
}
