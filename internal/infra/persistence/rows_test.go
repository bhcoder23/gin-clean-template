package persistence

import (
	"errors"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

var (
	errRows               = errors.New("rows error")
	errUnexpectedScanDest = errors.New("unexpected scan destination")
)

type fakeRows struct {
	index         int
	err           error
	tasks         []domain.Task
	notifications []domain.Notification
}

func (r *fakeRows) Close() {}

func (r *fakeRows) Err() error {
	return r.err
}

func (r *fakeRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *fakeRows) Next() bool {
	if len(r.tasks) > 0 {
		return r.index < len(r.tasks)
	}

	return r.index < len(r.notifications)
}

func (r *fakeRows) Scan(dest ...any) error {
	if len(r.tasks) > 0 {
		task := r.tasks[r.index]
		r.index++

		return assignScanDestinations(dest,
			task.ID,
			task.UserID,
			task.Title,
			task.Description,
			task.Status,
			task.CreatedAt,
			task.UpdatedAt,
		)
	}

	notification := r.notifications[r.index]
	r.index++

	return assignScanDestinations(dest,
		notification.ID,
		notification.UserID,
		notification.TaskID,
		notification.Type,
		notification.Title,
		notification.Body,
		notification.Read,
		notification.CreatedAt,
		notification.ReadAt,
	)
}

func assignScanDestinations(dest []any, values ...any) error { //nolint:cyclop,gocognit,gocyclo // fake pgx rows must cover several scan destination types.
	for i, value := range values {
		switch typed := dest[i].(type) {
		case *string:
			v, ok := value.(string)
			if !ok {
				return errUnexpectedScanDest
			}

			*typed = v
		case *domain.TaskStatus:
			v, ok := value.(domain.TaskStatus)
			if !ok {
				return errUnexpectedScanDest
			}

			*typed = v
		case *domain.NotificationType:
			v, ok := value.(domain.NotificationType)
			if !ok {
				return errUnexpectedScanDest
			}

			*typed = v
		case *bool:
			v, ok := value.(bool)
			if !ok {
				return errUnexpectedScanDest
			}

			*typed = v
		case *time.Time:
			v, ok := value.(time.Time)
			if !ok {
				return errUnexpectedScanDest
			}

			*typed = v
		case **time.Time:
			v, ok := value.(*time.Time)
			if !ok {
				return errUnexpectedScanDest
			}

			*typed = v
		default:
			return errUnexpectedScanDest
		}
	}

	return nil
}

func (r *fakeRows) Values() ([]any, error) {
	return nil, nil
}

func (r *fakeRows) RawValues() [][]byte {
	return nil
}

func (r *fakeRows) Conn() *pgx.Conn {
	return nil
}

var _ pgx.Rows = (*fakeRows)(nil)

func TestCollectTaskRowsChecksIteratorError(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	rows := &fakeRows{
		err: errRows,
		tasks: []domain.Task{
			{
				ID:        "task-1",
				UserID:    "user-1",
				Title:     "Task",
				Status:    domain.TaskStatusTodo,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	tasks, err := collectTaskRows(rows, 10)

	require.Nil(t, tasks)
	require.ErrorIs(t, err, errRows)
}

func TestCollectNotificationRowsChecksIteratorError(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	rows := &fakeRows{
		err: errRows,
		notifications: []domain.Notification{
			{
				ID:        "notification-1",
				UserID:    "user-1",
				TaskID:    "task-1",
				Type:      domain.NotificationTypeTaskCreated,
				Title:     "Task created",
				Body:      "Task created.",
				Read:      false,
				CreatedAt: now,
			},
		},
	}

	notifications, err := collectNotificationRows(rows, 10)

	require.Nil(t, notifications)
	require.ErrorIs(t, err, errRows)
}
