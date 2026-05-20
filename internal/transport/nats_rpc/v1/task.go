package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/nats-io/nats.go"
)

func (r *V1) createTask() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.CreateTaskReq

		if err = r.decodeAndValidate(data, &req); err != nil {
			return nil, err
		}

		task, err := r.tk.Create(ctx, userID, req.Title, req.Description)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - createTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(task), nil
	}
}

func (r *V1) getTask() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.GetTaskReq

		if err = r.decodeAndValidate(data, &req); err != nil {
			return nil, err
		}

		task, err := r.tk.Get(ctx, userID, req.ID)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - getTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(task), nil
	}
}

func (r *V1) listTasks() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.ListTasksReq

		if err = r.decodeAndValidate(data, &req); err != nil {
			return nil, err
		}

		var status *domain.TaskStatus

		if req.Status != "" {
			s := domain.TaskStatus(req.Status)
			status = &s
		}

		tasks, total, err := r.tk.List(ctx, userID, status, req.Query, req.Limit, req.Offset)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - listTasks")

			return nil, apperror.RPC(err)
		}

		return response.NewListTasksResp(tasks, total), nil
	}
}

func (r *V1) updateTask() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.UpdateTaskReq

		if err = r.decodeAndValidate(data, &req); err != nil {
			return nil, err
		}

		task, err := r.tk.Update(ctx, userID, req.ID, req.Title, req.Description)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - updateTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(task), nil
	}
}

func (r *V1) transitionTask() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.TransitionTaskReq

		if err = r.decodeAndValidate(data, &req); err != nil {
			return nil, err
		}

		task, err := r.tk.Transition(ctx, userID, req.ID, domain.TaskStatus(req.Status))
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - transitionTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(task), nil
	}
}

func (r *V1) deleteTask() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.DeleteTaskReq

		if err = r.decodeAndValidate(data, &req); err != nil {
			return nil, err
		}

		err = r.tk.Delete(ctx, userID, req.ID)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - deleteTask")

			return nil, apperror.RPC(err)
		}

		return response.DeleteStatus{Status: "deleted"}, nil
	}
}
