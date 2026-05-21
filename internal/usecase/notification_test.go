package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newNotificationUseCase(t *testing.T) (*notification.Usecase, *MockNotificationRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockRepo := NewMockNotificationRepo(ctrl)

	return notification.New(mockRepo), mockRepo
}

func TestNotificationList(t *testing.T) {
	t.Parallel()

	unreadOnly := true
	expected := []domain.Notification{
		{
			ID:        "notification-1",
			UserID:    "user-id-123",
			TaskID:    "task-id-123",
			Type:      domain.NotificationTypeTaskCreated,
			Title:     "Task created",
			Body:      "Task \"Ship the scaffold\" was created.",
			Read:      false,
			CreatedAt: time.Now().UTC(),
		},
	}

	uc, mockRepo := newNotificationUseCase(t)
	mockRepo.EXPECT().List(context.Background(), "user-id-123", appports.NotificationFilter{
		UnreadOnly: &unreadOnly,
		Limit:      uint64(10),
		Offset:     uint64(0),
	}).Return(expected, 1, nil)

	notifications, total, err := uc.List(context.Background(), "user-id-123", &unreadOnly, 10, 0)

	require.NoError(t, err)
	assert.Equal(t, expected, notifications)
	assert.Equal(t, 1, total)
}

func TestNotificationMarkRead(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newNotificationUseCase(t)
	existing := domain.Notification{
		ID:        "notification-1",
		UserID:    "user-id-123",
		TaskID:    "task-id-123",
		Type:      domain.NotificationTypeTaskCreated,
		Title:     "Task created",
		Body:      "Task \"Ship the scaffold\" was created.",
		Read:      false,
		CreatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "notification-1").Return(existing, nil)
	mockRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)

	updated, err := uc.MarkRead(context.Background(), "user-id-123", "notification-1")

	require.NoError(t, err)
	assert.True(t, updated.Read)
	require.NotNil(t, updated.ReadAt)
}

func TestNotificationMarkReadAlreadyRead(t *testing.T) {
	t.Parallel()

	readAt := time.Now().UTC()
	uc, mockRepo := newNotificationUseCase(t)
	existing := domain.Notification{
		ID:        "notification-1",
		UserID:    "user-id-123",
		TaskID:    "task-id-123",
		Type:      domain.NotificationTypeTaskCreated,
		Title:     "Task created",
		Body:      "Task \"Ship the scaffold\" was created.",
		Read:      true,
		ReadAt:    &readAt,
		CreatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "notification-1").Return(existing, nil)

	updated, err := uc.MarkRead(context.Background(), "user-id-123", "notification-1")

	require.NoError(t, err)
	assert.Equal(t, existing, updated)
}
