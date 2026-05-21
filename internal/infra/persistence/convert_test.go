package persistence

import (
	"math"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/infra/persistence/sqlc"
	"github.com/stretchr/testify/require"
)

func TestTaskToDomainConvertsNullableDescription(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	task := taskToDomain(&sqlc.Task{
		ID:        "task-1",
		UserID:    "user-1",
		Title:     "Task",
		Status:    string(domain.TaskStatusTodo),
		CreatedAt: now,
		UpdatedAt: now,
	})

	require.Empty(t, task.Description)
}

func TestNotificationToDomainPreservesReadAtPointer(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	readAt := now.Add(time.Minute)

	notification := notificationToDomain(&sqlc.Notification{
		ID:        "notification-1",
		UserID:    "user-1",
		TaskID:    "task-1",
		Type:      string(domain.NotificationTypeTaskCreated),
		Title:     "Task created",
		Body:      "Task created.",
		CreatedAt: now,
		ReadAt:    &readAt,
	})

	require.Equal(t, &readAt, notification.ReadAt)
}

func TestPaginationCountsRejectsInt64Overflow(t *testing.T) {
	t.Parallel()

	_, _, err := paginationCounts(uint64(math.MaxInt64)+1, 0)
	require.Error(t, err)
	require.ErrorContains(t, err, "limit")

	_, _, err = paginationCounts(10, uint64(math.MaxInt64)+1)
	require.Error(t, err)
	require.ErrorContains(t, err, "offset")
}
