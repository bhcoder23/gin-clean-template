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

// UseCase -.
type UseCase struct {
	repo             appports.TaskStore
	notificationRepo appports.NotificationStore
}

// New -.
func New(r appports.TaskStore, notificationRepo appports.NotificationStore) *UseCase {
	return &UseCase{
		repo:             r,
		notificationRepo: notificationRepo,
	}
}

// Create -.
func (uc *UseCase) Create(ctx context.Context, userID, title, description string) (domain.Task, error) {
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

	err := uc.repo.Store(ctx, &task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("TaskUseCase - Create - uc.repo.Store: %w", err)
	}

	if err = uc.storeNotification(ctx, domain.Notification{
		UserID: userID,
		TaskID: task.ID,
		Type:   domain.NotificationTypeTaskCreated,
		Title:  "Task created",
		Body:   fmt.Sprintf("Task %q was created with status %s.", task.Title, task.Status),
		Read:   false,
	}); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

// Get -.
func (uc *UseCase) Get(ctx context.Context, userID, taskID string) (domain.Task, error) {
	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("TaskUseCase - Get - uc.repo.GetByID: %w", err)
	}

	return task, nil
}

// List -.
func (uc *UseCase) List(ctx context.Context, userID string, status *domain.TaskStatus, query string, limit, offset int) ([]domain.Task, int, error) {
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
		return nil, 0, fmt.Errorf("TaskUseCase - List - uc.repo.List: %w", err)
	}

	return tasks, total, nil
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, userID, taskID, title, description string) (domain.Task, error) {
	now := time.Now().UTC()
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return domain.Task{}, domain.ErrTaskTitleRequired
	}

	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("TaskUseCase - Update - uc.repo.GetByID: %w", err)
	}

	if task.Status == domain.TaskStatusDone {
		return domain.Task{}, domain.ErrTaskCompleted
	}

	task.Title = title
	task.Description = description
	task.UpdatedAt = now

	err = uc.repo.Update(ctx, &task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("TaskUseCase - Update - uc.repo.Update: %w", err)
	}

	return task, nil
}

// Transition -.
func (uc *UseCase) Transition(ctx context.Context, userID, taskID string, newStatus domain.TaskStatus) (domain.Task, error) {
	now := time.Now().UTC()

	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("TaskUseCase - Transition - uc.repo.GetByID: %w", err)
	}

	err = task.Transition(newStatus)
	if err != nil {
		return domain.Task{}, err
	}

	task.UpdatedAt = now

	err = uc.repo.Update(ctx, &task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("TaskUseCase - Transition - uc.repo.Update: %w", err)
	}

	if err = uc.storeNotification(ctx, domain.Notification{
		UserID: userID,
		TaskID: task.ID,
		Type:   domain.NotificationTypeTaskStatusChanged,
		Title:  "Task status changed",
		Body:   fmt.Sprintf("Task %q moved to %s.", task.Title, task.Status),
		Read:   false,
	}); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

// Delete -.
func (uc *UseCase) Delete(ctx context.Context, userID, taskID string) error {
	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("TaskUseCase - Delete - uc.repo.GetByID: %w", err)
	}

	if task.Status == domain.TaskStatusDone {
		return domain.ErrTaskCompleted
	}

	err = uc.repo.Delete(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("TaskUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}

func (uc *UseCase) storeNotification(ctx context.Context, notification domain.Notification) error {
	notification.ID = uuid.New().String()
	notification.CreatedAt = time.Now().UTC()

	if err := uc.notificationRepo.Store(ctx, &notification); err != nil {
		return fmt.Errorf("TaskUseCase - storeNotification - uc.notificationRepo.Store: %w", err)
	}

	return nil
}
