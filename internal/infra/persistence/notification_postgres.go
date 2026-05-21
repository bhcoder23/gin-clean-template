package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/infra/persistence/sqlc"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// NotificationRepo -.
type NotificationRepo struct {
	queries *sqlc.Queries
}

// NewNotificationRepo -.
func NewNotificationRepo(pg *postgres.Postgres) *NotificationRepo {
	return NewNotificationRepoWithExecutor(pg.Pool)
}

// NewNotificationRepoWithExecutor creates a repository bound to a pool or transaction executor.
func NewNotificationRepoWithExecutor(executor postgres.Executor) *NotificationRepo {
	return &NotificationRepo{
		queries: sqlc.New(executor),
	}
}

// Store -.
func (r *NotificationRepo) Store(ctx context.Context, notification *domain.Notification) error {
	err := r.queries.CreateNotification(ctx, sqlc.CreateNotificationParams{
		ID:        notification.ID,
		UserID:    notification.UserID,
		TaskID:    notification.TaskID,
		Type:      string(notification.Type),
		Title:     notification.Title,
		Body:      notification.Body,
		Read:      notification.Read,
		CreatedAt: notification.CreatedAt,
		ReadAt:    notification.ReadAt,
	})
	if err != nil {
		return fmt.Errorf("NotificationRepo - Store - CreateNotification: %w", err)
	}

	return nil
}

// GetByID -.
func (r *NotificationRepo) GetByID(ctx context.Context, userID, notificationID string) (domain.Notification, error) {
	row, err := r.queries.GetNotificationByID(ctx, sqlc.GetNotificationByIDParams{
		ID:     notificationID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Notification{}, domain.ErrNotificationNotFound
		}

		return domain.Notification{}, fmt.Errorf("NotificationRepo - GetByID - GetNotificationByID: %w", err)
	}

	return notificationToDomain(&row), nil
}

// List -.
func (r *NotificationRepo) List(ctx context.Context, userID string, filter appports.NotificationFilter) ([]domain.Notification, int, error) {
	unreadOnly := filter.UnreadOnly != nil && *filter.UnreadOnly

	limitCount, offsetCount, err := paginationCounts(filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - paginationCounts: %w", err)
	}

	total, err := r.queries.CountNotifications(ctx, sqlc.CountNotificationsParams{
		UserID:     userID,
		UnreadOnly: unreadOnly,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - CountNotifications: %w", err)
	}

	rows, err := r.queries.ListNotifications(ctx, sqlc.ListNotificationsParams{
		UserID:      userID,
		UnreadOnly:  unreadOnly,
		OffsetCount: offsetCount,
		LimitCount:  limitCount,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - ListNotifications: %w", err)
	}

	return notificationsToDomain(rows), int(total), nil
}

// Update -.
func (r *NotificationRepo) Update(ctx context.Context, notification *domain.Notification) error {
	rowsAffected, err := r.queries.UpdateNotification(ctx, sqlc.UpdateNotificationParams{
		Read:   notification.Read,
		ReadAt: notification.ReadAt,
		ID:     notification.ID,
		UserID: notification.UserID,
	})
	if err != nil {
		return fmt.Errorf("NotificationRepo - Update - UpdateNotification: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotificationNotFound
	}

	return nil
}

func notificationToDomain(row *sqlc.Notification) domain.Notification {
	return domain.Notification{
		ID:        row.ID,
		UserID:    row.UserID,
		TaskID:    row.TaskID,
		Type:      domain.NotificationType(row.Type),
		Title:     row.Title,
		Body:      row.Body,
		Read:      row.Read,
		CreatedAt: row.CreatedAt,
		ReadAt:    row.ReadAt,
	}
}

func notificationsToDomain(rows []sqlc.Notification) []domain.Notification {
	notifications := make([]domain.Notification, 0, len(rows))
	for i := range rows {
		notifications = append(notifications, notificationToDomain(&rows[i]))
	}

	return notifications
}
