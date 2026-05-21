package response

import (
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewUserResp(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	user := domain.User{
		ID:        "user-id-123",
		Username:  "alice",
		Email:     "alice@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := NewUserResp(&user)

	require.Equal(t, "user-id-123", resp.ID)
	require.Equal(t, "alice", resp.Username)
	require.Equal(t, "alice@example.com", resp.Email)
	require.Equal(t, now, resp.CreatedAt)
	require.Equal(t, now, resp.UpdatedAt)
}

func TestNewTaskResp(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

	task := domain.Task{
		ID:          "task-id-123",
		UserID:      "user-id-123",
		Title:       "Ship scaffold",
		Description: "Keep boundaries clear",
		Status:      domain.TaskStatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	resp := NewTaskResp(&task)

	require.Equal(t, "task-id-123", resp.ID)
	require.Equal(t, "user-id-123", resp.UserID)
	require.Equal(t, "Ship scaffold", resp.Title)
	require.Equal(t, "Keep boundaries clear", resp.Description)
	require.Equal(t, string(domain.TaskStatusTodo), resp.Status)
	require.Equal(t, now, resp.CreatedAt)
	require.Equal(t, now, resp.UpdatedAt)
}

func TestNewNotificationResp(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)
	readAt := createdAt.Add(time.Hour)

	notification := domain.Notification{
		ID:        "notification-id-123",
		UserID:    "user-id-123",
		TaskID:    "task-id-123",
		Type:      domain.NotificationTypeTaskCreated,
		Title:     "Task created",
		Body:      "Task was created.",
		Read:      true,
		CreatedAt: createdAt,
		ReadAt:    &readAt,
	}

	resp := NewNotificationResp(&notification)

	require.Equal(t, "notification-id-123", resp.ID)
	require.Equal(t, "user-id-123", resp.UserID)
	require.Equal(t, "task-id-123", resp.TaskID)
	require.Equal(t, string(domain.NotificationTypeTaskCreated), resp.Type)
	require.Equal(t, "Task created", resp.Title)
	require.Equal(t, "Task was created.", resp.Body)
	require.True(t, resp.Read)
	require.Equal(t, createdAt, resp.CreatedAt)
	require.Equal(t, &readAt, resp.ReadAt)
}
