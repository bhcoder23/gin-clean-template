package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bhcoder23/gin-clean-template/internal/entity"
	"github.com/bhcoder23/gin-clean-template/internal/repo"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errRepoGeneric = errors.New("repository error")

func newTaskUseCase(t *testing.T) (*task.UseCase, *MockTaskRepo, *MockNotificationRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockRepo := NewMockTaskRepo(ctrl)
	notificationRepo := NewMockNotificationRepo(ctrl)
	useCase := task.New(mockRepo, notificationRepo)

	return useCase, mockRepo, notificationRepo
}

func TestTaskCreate(t *testing.T) {
	t.Parallel()

	t.Run("create success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, notificationRepo := newTaskUseCase(t)
		mockRepo.EXPECT().Store(context.Background(), gomock.Any()).Return(nil)
		notificationRepo.EXPECT().Store(context.Background(), gomock.Any()).Return(nil)

		t2, err := uc.Create(context.Background(), "user-id-123", "My Task", "Task description")

		require.NoError(t, err)
		assert.NotEmpty(t, t2.ID)
		assert.Equal(t, "My Task", t2.Title)
		assert.Equal(t, entity.TaskStatusTodo, t2.Status)
	})

	t.Run("create trims fields", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, notificationRepo := newTaskUseCase(t)
		mockRepo.EXPECT().Store(context.Background(), gomock.Any()).DoAndReturn(func(_ context.Context, task *entity.Task) error {
			assert.Equal(t, "My Task", task.Title)
			assert.Equal(t, "Task description", task.Description)

			return nil
		})
		notificationRepo.EXPECT().Store(context.Background(), gomock.Any()).Return(nil)

		created, err := uc.Create(context.Background(), "user-id-123", "  My Task  ", "  Task description  ")

		require.NoError(t, err)
		assert.Equal(t, "My Task", created.Title)
		assert.Equal(t, "Task description", created.Description)
	})

	t.Run("create rejects blank title after trim", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, notificationRepo := newTaskUseCase(t)
		mockRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Times(0)
		notificationRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Times(0)

		_, err := uc.Create(context.Background(), "user-id-123", "   ", "desc")

		require.ErrorIs(t, err, entity.ErrTaskTitleRequired)
	})
}

func TestTaskGet(t *testing.T) {
	t.Parallel()

	expectedTask := entity.Task{
		ID:     "task-id-123",
		UserID: "user-id-123",
		Title:  "My Task",
		Status: entity.TaskStatusTodo,
	}

	t.Run("get success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)
		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(expectedTask, nil)

		t2, err := uc.Get(context.Background(), "user-id-123", "task-id-123")

		require.NoError(t, err)
		assert.Equal(t, expectedTask, t2)
	})

	t.Run("get not found", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)
		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "missing-id").Return(entity.Task{}, entity.ErrTaskNotFound)

		_, err := uc.Get(context.Background(), "user-id-123", "missing-id")

		require.ErrorIs(t, err, entity.ErrTaskNotFound)
	})
}

func TestTaskList(t *testing.T) {
	t.Parallel()

	task1 := entity.Task{ID: "task-1", UserID: "user-id-123", Title: "Task 1", Status: entity.TaskStatusTodo}
	task2 := entity.Task{ID: "task-2", UserID: "user-id-123", Title: "Task 2", Status: entity.TaskStatusInProgress}

	t.Run("list success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)
		mockRepo.EXPECT().List(context.Background(), "user-id-123", gomock.Any()).Return([]entity.Task{task1, task2}, 2, nil)

		tasks, total, err := uc.List(context.Background(), "user-id-123", nil, "", 10, 0)

		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, tasks, 2)
	})

	t.Run("list defaults", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)
		mockRepo.EXPECT().List(context.Background(), "user-id-123", repo.TaskFilter{
			Status: nil,
			Limit:  uint64(10),
			Offset: uint64(0),
		}).Return([]entity.Task{task1, task2}, 2, nil)

		tasks, total, err := uc.List(context.Background(), "user-id-123", nil, "", 0, -1)

		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, tasks, 2)
	})

	t.Run("list trims query", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)
		mockRepo.EXPECT().List(context.Background(), "user-id-123", repo.TaskFilter{
			Status: nil,
			Query:  "urgent",
			Limit:  uint64(10),
			Offset: uint64(0),
		}).Return([]entity.Task{task1}, 1, nil)

		tasks, total, err := uc.List(context.Background(), "user-id-123", nil, "  urgent  ", 10, 0)

		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, tasks, 1)
	})
}

func TestTaskUpdate(t *testing.T) {
	t.Parallel()

	t.Run("update success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)

		existingTask := entity.Task{
			ID:     "task-id-123",
			UserID: "user-id-123",
			Title:  "Old Title",
			Status: entity.TaskStatusTodo,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(existingTask, nil)
		mockRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)

		updated, err := uc.Update(context.Background(), "user-id-123", "task-id-123", "New Title", "New description")

		require.NoError(t, err)
		assert.Equal(t, "New Title", updated.Title)
	})

	t.Run("update trims fields", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)

		existingTask := entity.Task{
			ID:          "task-id-123",
			UserID:      "user-id-123",
			Title:       "Old Title",
			Description: "Old description",
			Status:      entity.TaskStatusTodo,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(existingTask, nil)
		mockRepo.EXPECT().Update(context.Background(), gomock.Any()).DoAndReturn(func(_ context.Context, task *entity.Task) error {
			assert.Equal(t, "New Title", task.Title)
			assert.Equal(t, "New description", task.Description)

			return nil
		})

		updated, err := uc.Update(context.Background(), "user-id-123", "task-id-123", "  New Title  ", "  New description  ")

		require.NoError(t, err)
		assert.Equal(t, "New Title", updated.Title)
		assert.Equal(t, "New description", updated.Description)
	})

	t.Run("update rejects done task", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)

		existingTask := entity.Task{
			ID:     "task-id-123",
			UserID: "user-id-123",
			Title:  "Done Task",
			Status: entity.TaskStatusDone,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(existingTask, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		_, err := uc.Update(context.Background(), "user-id-123", "task-id-123", "New Title", "desc")

		require.ErrorIs(t, err, entity.ErrTaskCompleted)
	})
}

func TestTaskTransition(t *testing.T) {
	t.Parallel()

	t.Run("transition valid", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, notificationRepo := newTaskUseCase(t)

		todoTask := entity.Task{
			ID:     "task-id-123",
			UserID: "user-id-123",
			Title:  "My Task",
			Status: entity.TaskStatusTodo,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(todoTask, nil)
		mockRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
		notificationRepo.EXPECT().Store(context.Background(), gomock.Any()).Return(nil)

		updated, err := uc.Transition(context.Background(), "user-id-123", "task-id-123", entity.TaskStatusInProgress)

		require.NoError(t, err)
		assert.Equal(t, entity.TaskStatusInProgress, updated.Status)
	})

	t.Run("transition invalid", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)

		doneTask := entity.Task{
			ID:     "task-id-456",
			UserID: "user-id-123",
			Title:  "Done Task",
			Status: entity.TaskStatusDone,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-456").Return(doneTask, nil)

		_, err := uc.Transition(context.Background(), "user-id-123", "task-id-456", entity.TaskStatusTodo)

		require.ErrorIs(t, err, entity.ErrInvalidTransition)
	})

	t.Run("transition done task cannot move again", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, notificationRepo := newTaskUseCase(t)

		doneTask := entity.Task{
			ID:     "task-id-456",
			UserID: "user-id-123",
			Title:  "Done Task",
			Status: entity.TaskStatusDone,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-456").Return(doneTask, nil)
		notificationRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Times(0)

		_, err := uc.Transition(context.Background(), "user-id-123", "task-id-456", entity.TaskStatusTodo)

		require.ErrorIs(t, err, entity.ErrInvalidTransition)
	})
}

func TestTaskDelete(t *testing.T) {
	t.Parallel()

	t.Run("delete success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)
		activeTask := entity.Task{
			ID:     "task-id-123",
			UserID: "user-id-123",
			Title:  "Task",
			Status: entity.TaskStatusTodo,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(activeTask, nil)
		mockRepo.EXPECT().Delete(context.Background(), "user-id-123", "task-id-123").Return(nil)

		err := uc.Delete(context.Background(), "user-id-123", "task-id-123")

		require.NoError(t, err)
	})

	t.Run("delete not found", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)
		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "missing-id").Return(entity.Task{}, entity.ErrTaskNotFound)

		err := uc.Delete(context.Background(), "user-id-123", "missing-id")

		require.ErrorIs(t, err, entity.ErrTaskNotFound)
	})

	t.Run("delete rejects done task", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newTaskUseCase(t)

		doneTask := entity.Task{
			ID:     "task-id-123",
			UserID: "user-id-123",
			Title:  "Done Task",
			Status: entity.TaskStatusDone,
		}

		mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(doneTask, nil)
		mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		err := uc.Delete(context.Background(), "user-id-123", "task-id-123")

		require.ErrorIs(t, err, entity.ErrTaskCompleted)
	})
}

func TestTaskCreate_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)

	mockRepo.EXPECT().Store(context.Background(), gomock.Any()).Return(errRepoGeneric)

	_, err := uc.Create(context.Background(), "user-id-123", "title", "desc")

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskGet_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)

	mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-999").Return(entity.Task{}, errRepoGeneric)

	_, err := uc.Get(context.Background(), "user-id-123", "task-id-999")

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskUpdate_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)

	existing := entity.Task{
		ID:     "task-id-123",
		UserID: "user-id-123",
		Title:  "Old Title",
		Status: entity.TaskStatusTodo,
	}

	mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(existing, nil)
	mockRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(errRepoGeneric)

	_, err := uc.Update(context.Background(), "user-id-123", "task-id-123", "New Title", "desc")

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskUpdate_NotFound(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)

	mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "missing-id").Return(entity.Task{}, entity.ErrTaskNotFound)

	_, err := uc.Update(context.Background(), "user-id-123", "missing-id", "title", "desc")

	require.Error(t, err)
	require.ErrorIs(t, err, entity.ErrTaskNotFound)
}

func TestTaskTransition_UpdateError(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)

	todoTask := entity.Task{
		ID:     "task-id-123",
		UserID: "user-id-123",
		Title:  "My Task",
		Status: entity.TaskStatusTodo,
	}

	mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(todoTask, nil)
	mockRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(errRepoGeneric)

	_, err := uc.Transition(context.Background(), "user-id-123", "task-id-123", entity.TaskStatusInProgress)

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskDelete_GenericError(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)
	activeTask := entity.Task{
		ID:     "task-id-123",
		UserID: "user-id-123",
		Title:  "Task",
		Status: entity.TaskStatusTodo,
	}

	mockRepo.EXPECT().GetByID(context.Background(), "user-id-123", "task-id-123").Return(activeTask, nil)
	mockRepo.EXPECT().Delete(context.Background(), "user-id-123", "task-id-123").Return(errRepoGeneric)

	err := uc.Delete(context.Background(), "user-id-123", "task-id-123")

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskList_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)

	mockRepo.EXPECT().
		List(context.Background(), "user-id-123", repo.TaskFilter{Limit: uint64(10), Offset: uint64(0)}).
		Return(nil, 0, errRepoGeneric)

	_, _, err := uc.List(context.Background(), "user-id-123", nil, "", 10, 0)

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskTransition_NotFound(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newTaskUseCase(t)

	mockRepo.EXPECT().
		GetByID(context.Background(), "user-id-123", "task-id-123").
		Return(entity.Task{}, entity.ErrTaskNotFound)

	_, err := uc.Transition(context.Background(), "user-id-123", "task-id-123", entity.TaskStatusInProgress)

	require.Error(t, err)
	require.ErrorIs(t, err, entity.ErrTaskNotFound)
}
