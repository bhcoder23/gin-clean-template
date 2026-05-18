package v1

import (
	"context"

	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/errlog"
	"github.com/bhcoder23/gin-clean-template/internal/transport/errmap"
	grpcmw "github.com/bhcoder23/gin-clean-template/internal/transport/grpc/middleware"
	"github.com/bhcoder23/gin-clean-template/internal/transport/grpc/v1/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateTask -.
func (c *TaskController) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.TaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	task, err := c.tk.Create(ctx, userID, req.GetTitle(), req.GetDescription())
	if err != nil {
		errlog.Log(c.l, err, "grpc - v1 - CreateTask")
		return nil, errmap.GRPC(err)
	}

	return response.NewTaskResponse(&task), nil
}

// GetTask -.
func (c *TaskController) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.TaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	task, err := c.tk.Get(ctx, userID, req.GetId())
	if err != nil {
		errlog.Log(c.l, err, "grpc - v1 - GetTask")
		return nil, errmap.GRPC(err)
	}

	return response.NewTaskResponse(&task), nil
}

// ListTasks -.
func (c *TaskController) ListTasks(ctx context.Context, req *v1.ListTasksRequest) (*v1.ListTasksResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	var statusFilter *domain.TaskStatus

	if req.GetStatus() != "" {
		s := domain.TaskStatus(req.GetStatus())
		if !s.Valid() {
			return nil, status.Error(codes.InvalidArgument, "invalid task status")
		}

		statusFilter = &s
	}

	tasks, total, err := c.tk.List(ctx, userID, statusFilter, req.GetQuery(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		errlog.Log(c.l, err, "grpc - v1 - ListTasks")

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return response.NewListTasksResponse(tasks, total), nil
}

// UpdateTask -.
func (c *TaskController) UpdateTask(ctx context.Context, req *v1.UpdateTaskRequest) (*v1.TaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	task, err := c.tk.Update(ctx, userID, req.GetId(), req.GetTitle(), req.GetDescription())
	if err != nil {
		errlog.Log(c.l, err, "grpc - v1 - UpdateTask")
		return nil, errmap.GRPC(err)
	}

	return response.NewTaskResponse(&task), nil
}

// TransitionTask -.
func (c *TaskController) TransitionTask(ctx context.Context, req *v1.TransitionTaskRequest) (*v1.TaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	task, err := c.tk.Transition(ctx, userID, req.GetId(), domain.TaskStatus(req.GetStatus()))
	if err != nil {
		errlog.Log(c.l, err, "grpc - v1 - TransitionTask")
		return nil, errmap.GRPC(err)
	}

	return response.NewTaskResponse(&task), nil
}

// DeleteTask -.
func (c *TaskController) DeleteTask(ctx context.Context, req *v1.DeleteTaskRequest) (*v1.DeleteTaskResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err := c.tk.Delete(ctx, userID, req.GetId())
	if err != nil {
		errlog.Log(c.l, err, "grpc - v1 - DeleteTask")
		return nil, errmap.GRPC(err)
	}

	return &v1.DeleteTaskResponse{}, nil
}
