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

var errRows = errors.New("rows error")

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
		if r.index < len(r.tasks) {
			return true
		}

		return false
	}

	if r.index < len(r.notifications) {
		return true
	}

	return false
}

func (r *fakeRows) Scan(dest ...any) error {
	if len(r.tasks) > 0 {
		task := r.tasks[r.index]
		r.index++

		*(dest[0].(*string)) = task.ID
		*(dest[1].(*string)) = task.UserID
		*(dest[2].(*string)) = task.Title
		*(dest[3].(*string)) = task.Description
		*(dest[4].(*domain.TaskStatus)) = task.Status
		*(dest[5].(*time.Time)) = task.CreatedAt
		*(dest[6].(*time.Time)) = task.UpdatedAt

		return nil
	}

	notification := r.notifications[r.index]
	r.index++

	*(dest[0].(*string)) = notification.ID
	*(dest[1].(*string)) = notification.UserID
	*(dest[2].(*string)) = notification.TaskID
	*(dest[3].(*domain.NotificationType)) = notification.Type
	*(dest[4].(*string)) = notification.Title
	*(dest[5].(*string)) = notification.Body
	*(dest[6].(*bool)) = notification.Read
	*(dest[7].(*time.Time)) = notification.CreatedAt
	*(dest[8].(**time.Time)) = notification.ReadAt

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
