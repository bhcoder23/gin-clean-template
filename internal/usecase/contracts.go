// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

//go:generate go tool mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// NotificationRepo persists notification aggregates.
	NotificationRepo interface {
		Store(ctx context.Context, notification *domain.Notification) error
		GetByID(ctx context.Context, userID, notificationID string) (domain.Notification, error)
		List(ctx context.Context, userID string, filter NotificationFilter) ([]domain.Notification, int, error)
		Update(ctx context.Context, notification *domain.Notification) error
	}

	// UserRepo persists user aggregates.
	UserRepo interface {
		Store(ctx context.Context, user *domain.User) error
		GetByID(ctx context.Context, id string) (domain.User, error)
		GetByEmail(ctx context.Context, email string) (domain.User, error)
	}

	// TaskRepo persists task aggregates.
	TaskRepo interface {
		Store(ctx context.Context, task *domain.Task) error
		GetByID(ctx context.Context, userID, taskID string) (domain.Task, error)
		List(ctx context.Context, userID string, filter TaskFilter) ([]domain.Task, int, error)
		Update(ctx context.Context, task *domain.Task) error
		Delete(ctx context.Context, userID, taskID string) error
	}

	// OutboxStore stores integration events in the same transaction as business data.
	OutboxStore interface {
		Add(ctx context.Context, event *OutboxEvent) error
	}

	// RepoProvider exposes repositories bound to the same persistence executor.
	RepoProvider interface {
		Users() UserRepo
		Tasks() TaskRepo
		Notifications() NotificationRepo
		Outbox() OutboxStore
	}

	// Transactor runs multi-repository use cases in one transaction.
	Transactor interface {
		WithinTx(ctx context.Context, fn func(context.Context, RepoProvider) error) error
	}

	// TaskFilter -.
	TaskFilter struct {
		Status *domain.TaskStatus
		Query  string
		Limit  uint64
		Offset uint64
	}

	// NotificationFilter -.
	NotificationFilter struct {
		UnreadOnly *bool
		Limit      uint64
		Offset     uint64
	}

	// OutboxEvent describes a business integration event before relay metadata is added.
	OutboxEvent struct {
		ID            string
		AggregateType string
		AggregateID   string
		EventType     string
		Payload       []byte
		Headers       map[string]string
		AvailableAt   time.Time
	}

	// Notification -.
	Notification interface {
		List(ctx context.Context, userID string, unreadOnly *bool, limit, offset int) ([]domain.Notification, int, error)
		MarkRead(ctx context.Context, userID, notificationID string) (domain.Notification, error)
	}

	// User -.
	User interface {
		Register(ctx context.Context, username, email, password string) (domain.User, error)
		Login(ctx context.Context, email, password string) (string, error)
		GetUser(ctx context.Context, userID string) (domain.User, error)
	}

	// Task -.
	Task interface {
		Create(ctx context.Context, userID, title, description string) (domain.Task, error)
		Get(ctx context.Context, userID, taskID string) (domain.Task, error)
		List(ctx context.Context, userID string, status *domain.TaskStatus, query string, limit, offset int) ([]domain.Task, int, error)
		Update(ctx context.Context, userID, taskID, title, description string) (domain.Task, error)
		Transition(ctx context.Context, userID, taskID string, newStatus domain.TaskStatus) (domain.Task, error)
		Delete(ctx context.Context, userID, taskID string) error
	}
)
