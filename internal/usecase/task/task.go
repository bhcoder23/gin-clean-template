package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/google/uuid"
)

// Usecase coordinates task application workflows.
type Usecase struct {
	repo             appports.TaskRepo
	notificationRepo appports.NotificationRepo
	transactor       appports.Transactor
}

// New -.
func New(r appports.TaskRepo, notificationRepo appports.NotificationRepo, transactors ...appports.Transactor) *Usecase {
	var transactor appports.Transactor
	if len(transactors) > 0 {
		transactor = transactors[0]
	}

	return &Usecase{
		repo:             r,
		notificationRepo: notificationRepo,
		transactor:       transactor,
	}
}

// Create -.
func (uc *Usecase) Create(ctx context.Context, userID, title, description string) (domain.Task, error) {
	now := time.Now().UTC()
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return domain.Task{}, domain.ErrTaskTitleRequired
	}

	task := domain.Task{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      domain.TaskStatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := uc.withRepos(ctx, func(txCtx context.Context, taskRepo appports.TaskRepo, notificationRepo appports.NotificationRepo) error {
		if err := taskRepo.Store(txCtx, &task); err != nil {
			return fmt.Errorf("task.Usecase - Create - taskRepo.Store: %w", err)
		}

		if err := uc.storeNotification(txCtx, notificationRepo, &domain.Notification{
			UserID: userID,
			TaskID: task.ID,
			Type:   domain.NotificationTypeTaskCreated,
			Title:  "Task created",
			Body:   fmt.Sprintf("Task %q was created with status %s.", task.Title, task.Status),
			Read:   false,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

// Get -.
func (uc *Usecase) Get(ctx context.Context, userID, taskID string) (domain.Task, error) {
	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("task.Usecase - Get - uc.repo.GetByID: %w", err)
	}

	return task, nil
}

// List -.
func (uc *Usecase) List(ctx context.Context, userID string, status *domain.TaskStatus, query string, limit, offset int) ([]domain.Task, int, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	tasks, total, err := uc.repo.List(ctx, userID, appports.TaskFilter{
		Status: status,
		Query:  strings.TrimSpace(query),
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("task.Usecase - List - uc.repo.List: %w", err)
	}

	return tasks, total, nil
}

// Update -.
func (uc *Usecase) Update(ctx context.Context, userID, taskID, title, description string) (domain.Task, error) {
	now := time.Now().UTC()
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return domain.Task{}, domain.ErrTaskTitleRequired
	}

	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("task.Usecase - Update - uc.repo.GetByID: %w", err)
	}

	if task.Status == domain.TaskStatusDone {
		return domain.Task{}, domain.ErrTaskCompleted
	}

	task.Title = title
	task.Description = description
	task.UpdatedAt = now

	err = uc.repo.Update(ctx, &task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("task.Usecase - Update - uc.repo.Update: %w", err)
	}

	return task, nil
}

// Transition -.
func (uc *Usecase) Transition(ctx context.Context, userID, taskID string, newStatus domain.TaskStatus) (domain.Task, error) {
	now := time.Now().UTC()

	var task domain.Task

	err := uc.withRepos(ctx, func(txCtx context.Context, taskRepo appports.TaskRepo, notificationRepo appports.NotificationRepo) error {
		var err error

		task, err = taskRepo.GetByID(txCtx, userID, taskID)
		if err != nil {
			return fmt.Errorf("task.Usecase - Transition - taskRepo.GetByID: %w", err)
		}

		err = task.Transition(newStatus)
		if err != nil {
			return err
		}

		task.UpdatedAt = now

		err = taskRepo.Update(txCtx, &task)
		if err != nil {
			return fmt.Errorf("task.Usecase - Transition - taskRepo.Update: %w", err)
		}

		if err := uc.storeNotification(txCtx, notificationRepo, &domain.Notification{
			UserID: userID,
			TaskID: task.ID,
			Type:   domain.NotificationTypeTaskStatusChanged,
			Title:  "Task status changed",
			Body:   fmt.Sprintf("Task %q moved to %s.", task.Title, task.Status),
			Read:   false,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

// Delete -.
func (uc *Usecase) Delete(ctx context.Context, userID, taskID string) error {
	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("task.Usecase - Delete - uc.repo.GetByID: %w", err)
	}

	if task.Status == domain.TaskStatusDone {
		return domain.ErrTaskCompleted
	}

	err = uc.repo.Delete(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("task.Usecase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}

func (uc *Usecase) withRepos(
	ctx context.Context,
	fn func(context.Context, appports.TaskRepo, appports.NotificationRepo) error,
) error {
	if uc.transactor == nil {
		return fn(ctx, uc.repo, uc.notificationRepo)
	}

	return uc.transactor.WithinTx(ctx, func(txCtx context.Context, repos appports.RepoProvider) error {
		return fn(txCtx, repos.Tasks(), repos.Notifications())
	})
}

func (uc *Usecase) storeNotification(ctx context.Context, notificationRepo appports.NotificationRepo, notification *domain.Notification) error {
	notification.ID = uuid.New().String()
	notification.CreatedAt = time.Now().UTC()

	if err := notificationRepo.Store(ctx, notification); err != nil {
		return fmt.Errorf("task.Usecase - storeNotification - notificationRepo.Store: %w", err)
	}

	return nil
}
