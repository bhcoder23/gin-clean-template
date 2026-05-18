// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/entity"
)

//go:generate go tool mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// NotificationRepo -.
	NotificationRepo interface {
		Store(ctx context.Context, notification *entity.Notification) error
		GetByID(ctx context.Context, userID, notificationID string) (entity.Notification, error)
		List(ctx context.Context, userID string, filter NotificationFilter) ([]entity.Notification, int, error)
		Update(ctx context.Context, notification *entity.Notification) error
	}

	// UserRepo -.
	UserRepo interface {
		Store(ctx context.Context, user *entity.User) error
		GetByID(ctx context.Context, id string) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
	}

	// TaskRepo -.
	TaskRepo interface {
		Store(ctx context.Context, task *entity.Task) error
		GetByID(ctx context.Context, userID, taskID string) (entity.Task, error)
		List(ctx context.Context, userID string, filter TaskFilter) ([]entity.Task, int, error)
		Update(ctx context.Context, task *entity.Task) error
		Delete(ctx context.Context, userID, taskID string) error
	}

	// TaskFilter -.
	TaskFilter struct {
		Status *entity.TaskStatus
		Limit  uint64
		Offset uint64
	}

	// NotificationFilter -.
	NotificationFilter struct {
		UnreadOnly *bool
		Limit      uint64
		Offset     uint64
	}
)
