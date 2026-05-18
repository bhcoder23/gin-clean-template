package persistence

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// NotificationRepo -.
type NotificationRepo struct {
	*postgres.Postgres
}

func collectNotificationRows(rows pgx.Rows, limit uint64) ([]domain.Notification, error) {
	notifications := make([]domain.Notification, 0, limit)

	for rows.Next() {
		var n domain.Notification

		err := rows.Scan(&n.ID, &n.UserID, &n.TaskID, &n.Type, &n.Title, &n.Body, &n.Read, &n.CreatedAt, &n.ReadAt)
		if err != nil {
			return nil, fmt.Errorf("collectNotificationRows - rows.Scan: %w", err)
		}

		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("collectNotificationRows - rows.Err: %w", err)
	}

	return notifications, nil
}

// NewNotificationRepo -.
func NewNotificationRepo(pg *postgres.Postgres) *NotificationRepo {
	return &NotificationRepo{pg}
}

// Store -.
func (r *NotificationRepo) Store(ctx context.Context, notification *domain.Notification) error {
	sql, args, err := r.Builder.
		Insert("notifications").
		Columns("id, user_id, task_id, type, title, body, read, created_at, read_at").
		Values(notification.ID, notification.UserID, notification.TaskID, notification.Type, notification.Title, notification.Body, notification.Read, notification.CreatedAt, notification.ReadAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - Store - r.Builder: %w", err)
	}

	if _, err = r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("NotificationRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetByID -.
func (r *NotificationRepo) GetByID(ctx context.Context, userID, notificationID string) (domain.Notification, error) {
	sql, args, err := r.Builder.
		Select("id, user_id, task_id, type, title, body, read, created_at, read_at").
		From("notifications").
		Where(sq.Eq{"id": notificationID, "user_id": userID}).
		ToSql()
	if err != nil {
		return domain.Notification{}, fmt.Errorf("NotificationRepo - GetByID - r.Builder: %w", err)
	}

	var notification domain.Notification

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&notification.ID, &notification.UserID, &notification.TaskID, &notification.Type, &notification.Title, &notification.Body, &notification.Read, &notification.CreatedAt, &notification.ReadAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Notification{}, domain.ErrNotificationNotFound
		}

		return domain.Notification{}, fmt.Errorf("NotificationRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return notification, nil
}

// List -.
func (r *NotificationRepo) List(ctx context.Context, userID string, filter appports.NotificationFilter) ([]domain.Notification, int, error) {
	countBuilder := r.Builder.
		Select("COUNT(*)").
		From("notifications").
		Where(sq.Eq{"user_id": userID})

	if filter.UnreadOnly != nil && *filter.UnreadOnly {
		countBuilder = countBuilder.Where(sq.Eq{"read": false})
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - countBuilder: %w", err)
	}

	var total int
	if err = r.Pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - count query: %w", err)
	}

	dataBuilder := r.Builder.
		Select("id, user_id, task_id, type, title, body, read, created_at, read_at").
		From("notifications").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset)

	if filter.UnreadOnly != nil && *filter.UnreadOnly {
		dataBuilder = dataBuilder.Where(sq.Eq{"read": false})
	}

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - dataBuilder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	notifications, err := collectNotificationRows(rows, filter.Limit)
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - collectNotificationRows: %w", err)
	}

	return notifications, total, nil
}

// Update -.
func (r *NotificationRepo) Update(ctx context.Context, notification *domain.Notification) error {
	sql, args, err := r.Builder.
		Update("notifications").
		Set("read", notification.Read).
		Set("read_at", notification.ReadAt).
		Where(sq.Eq{"id": notification.ID, "user_id": notification.UserID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - Update - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - Update - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotificationNotFound
	}

	return nil
}
