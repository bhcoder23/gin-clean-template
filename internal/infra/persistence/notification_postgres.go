package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// NotificationRepo -.
type NotificationRepo struct {
	builder  sq.StatementBuilderType
	executor postgres.Executor
}

type notificationRow struct {
	id        string
	userID    string
	taskID    string
	typ       domain.NotificationType
	title     string
	body      string
	read      bool
	createdAt time.Time
	readAt    *time.Time
}

func newNotificationRow(notification domain.Notification) notificationRow {
	return notificationRow{
		id:        notification.ID,
		userID:    notification.UserID,
		taskID:    notification.TaskID,
		typ:       notification.Type,
		title:     notification.Title,
		body:      notification.Body,
		read:      notification.Read,
		createdAt: notification.CreatedAt,
		readAt:    notification.ReadAt,
	}
}

func (r notificationRow) toDomain() domain.Notification {
	return domain.Notification{
		ID:        r.id,
		UserID:    r.userID,
		TaskID:    r.taskID,
		Type:      r.typ,
		Title:     r.title,
		Body:      r.body,
		Read:      r.read,
		CreatedAt: r.createdAt,
		ReadAt:    r.readAt,
	}
}

func collectNotificationRows(rows pgx.Rows, limit uint64) ([]domain.Notification, error) {
	notifications := make([]domain.Notification, 0, limit)

	for rows.Next() {
		var row notificationRow

		err := rows.Scan(&row.id, &row.userID, &row.taskID, &row.typ, &row.title, &row.body, &row.read, &row.createdAt, &row.readAt)
		if err != nil {
			return nil, fmt.Errorf("collectNotificationRows - rows.Scan: %w", err)
		}

		notifications = append(notifications, row.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("collectNotificationRows - rows.Err: %w", err)
	}

	return notifications, nil
}

// NewNotificationRepo -.
func NewNotificationRepo(pg *postgres.Postgres) *NotificationRepo {
	return NewNotificationRepoWithExecutor(pg.Builder, pg.Pool)
}

// NewNotificationRepoWithExecutor creates a repository bound to a pool or transaction executor.
func NewNotificationRepoWithExecutor(builder sq.StatementBuilderType, executor postgres.Executor) *NotificationRepo {
	return &NotificationRepo{
		builder:  builder,
		executor: executor,
	}
}

// Store -.
func (r *NotificationRepo) Store(ctx context.Context, notification *domain.Notification) error {
	row := newNotificationRow(*notification)

	sql, args, err := r.builder.
		Insert("notifications").
		Columns("id, user_id, task_id, type, title, body, read, created_at, read_at").
		Values(row.id, row.userID, row.taskID, row.typ, row.title, row.body, row.read, row.createdAt, row.readAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - Store - r.Builder: %w", err)
	}

	if _, err = r.executor.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("NotificationRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetByID -.
func (r *NotificationRepo) GetByID(ctx context.Context, userID, notificationID string) (domain.Notification, error) {
	sql, args, err := r.builder.
		Select("id, user_id, task_id, type, title, body, read, created_at, read_at").
		From("notifications").
		Where(sq.Eq{"id": notificationID, "user_id": userID}).
		ToSql()
	if err != nil {
		return domain.Notification{}, fmt.Errorf("NotificationRepo - GetByID - r.Builder: %w", err)
	}

	var row notificationRow

	err = r.executor.QueryRow(ctx, sql, args...).
		Scan(&row.id, &row.userID, &row.taskID, &row.typ, &row.title, &row.body, &row.read, &row.createdAt, &row.readAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Notification{}, domain.ErrNotificationNotFound
		}

		return domain.Notification{}, fmt.Errorf("NotificationRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return row.toDomain(), nil
}

// List -.
func (r *NotificationRepo) List(ctx context.Context, userID string, filter appports.NotificationFilter) ([]domain.Notification, int, error) {
	countBuilder := r.builder.
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
	if err = r.executor.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo - List - count query: %w", err)
	}

	dataBuilder := r.builder.
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

	rows, err := r.executor.Query(ctx, dataSQL, dataArgs...)
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
	row := newNotificationRow(*notification)

	sql, args, err := r.builder.
		Update("notifications").
		Set("read", row.read).
		Set("read_at", row.readAt).
		Where(sq.Eq{"id": row.id, "user_id": row.userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("NotificationRepo - Update - r.Builder: %w", err)
	}

	result, err := r.executor.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepo - Update - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotificationNotFound
	}

	return nil
}
