package v1

import (
	"context"

	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	grpcmw "github.com/bhcoder23/gin-clean-template/internal/transport/grpc/middleware"
	"github.com/bhcoder23/gin-clean-template/internal/transport/grpc/v1/response"
)

// CreateTask -.
func (c *TaskController) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.TaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, apperror.GRPC(apperror.ErrUnauthorized)
	}

	task, err := c.tk.Create(ctx, userID, req.GetTitle(), req.GetDescription())
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - CreateTask")

		return nil, apperror.GRPC(err)
	}

	return response.NewTaskResponse(&task), nil
}

// GetTask -.
func (c *TaskController) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.TaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, apperror.GRPC(apperror.ErrUnauthorized)
	}

	task, err := c.tk.Get(ctx, userID, req.GetId())
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - GetTask")

		return nil, apperror.GRPC(err)
	}

	return response.NewTaskResponse(&task), nil
}

// ListTasks -.
func (c *TaskController) ListTasks(ctx context.Context, req *v1.ListTasksRequest) (*v1.ListTasksResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, apperror.GRPC(apperror.ErrUnauthorized)
	}

	var statusFilter *domain.TaskStatus

	if req.GetStatus() != "" {
		s := domain.TaskStatus(req.GetStatus())
		if !s.Valid() {
			return nil, apperror.GRPCWithMessage(apperror.ErrInvalidRequest, "invalid task status")
		}

		statusFilter = &s
	}

	tasks, total, err := c.tk.List(ctx, userID, statusFilter, req.GetQuery(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - ListTasks")

		return nil, apperror.GRPC(err)
	}

	return response.NewListTasksResponse(tasks, total), nil
}

// UpdateTask -.
func (c *TaskController) UpdateTask(ctx context.Context, req *v1.UpdateTaskRequest) (*v1.TaskResponse, error) {
	return c.writeTaskResponse(ctx, "grpc - v1 - UpdateTask", func(userID string) (domain.Task, error) {
		return c.tk.Update(ctx, userID, req.GetId(), req.GetTitle(), req.GetDescription())
	})
}

// TransitionTask -.
func (c *TaskController) TransitionTask(ctx context.Context, req *v1.TransitionTaskRequest) (*v1.TaskResponse, error) {
	return c.writeTaskResponse(ctx, "grpc - v1 - TransitionTask", func(userID string) (domain.Task, error) {
		return c.tk.Transition(ctx, userID, req.GetId(), domain.TaskStatus(req.GetStatus()))
	})
}

func (c *TaskController) writeTaskResponse(ctx context.Context, logMessage string, write func(string) (domain.Task, error)) (*v1.TaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, apperror.GRPC(apperror.ErrUnauthorized)
	}

	task, err := write(userID)
	if err != nil {
		apperror.Log(c.l, err, logMessage)

		return nil, apperror.GRPC(err)
	}

	return response.NewTaskResponse(&task), nil
}

// DeleteTask -.
func (c *TaskController) DeleteTask(ctx context.Context, req *v1.DeleteTaskRequest) (*v1.DeleteTaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, apperror.GRPC(apperror.ErrUnauthorized)
	}

	err := c.tk.Delete(ctx, userID, req.GetId())
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - DeleteTask")

		return nil, apperror.GRPC(err)
	}

	return &v1.DeleteTaskResponse{}, nil
}
