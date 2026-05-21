package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) createTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.CreateTaskReq

		userID, err := r.bindAuthenticatedRequest(d, &req)
		if err != nil {
			return nil, err
		}

		task, err := r.tk.Create(ctx, userID, req.Title, req.Description)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - createTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(&task), nil
	}
}

//nolint:dupl // RPC handlers stay explicit per route; shared binding/error mapping is already factored.
func (r *V1) getTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.GetTaskReq

		userID, err := r.bindAuthenticatedRequest(d, &req)
		if err != nil {
			return nil, err
		}

		task, err := r.tk.Get(ctx, userID, req.ID)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - getTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(&task), nil
	}
}

func (r *V1) listTasks() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.ListTasksReq

		userID, err := r.bindAuthenticatedRequest(d, &req)
		if err != nil {
			return nil, err
		}

		var status *domain.TaskStatus

		if req.Status != "" {
			s := domain.TaskStatus(req.Status)
			status = &s
		}

		tasks, total, err := r.tk.List(ctx, userID, status, req.Query, req.Limit, req.Offset)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - listTasks")

			return nil, apperror.RPC(err)
		}

		return response.NewListTasksResp(tasks, total), nil
	}
}

func (r *V1) updateTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.UpdateTaskReq

		userID, err := r.bindAuthenticatedRequest(d, &req)
		if err != nil {
			return nil, err
		}

		task, err := r.tk.Update(ctx, userID, req.ID, req.Title, req.Description)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - updateTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(&task), nil
	}
}

func (r *V1) transitionTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.TransitionTaskReq

		userID, err := r.bindAuthenticatedRequest(d, &req)
		if err != nil {
			return nil, err
		}

		task, err := r.tk.Transition(ctx, userID, req.ID, domain.TaskStatus(req.Status))
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - transitionTask")

			return nil, apperror.RPC(err)
		}

		return response.NewTaskResp(&task), nil
	}
}

func (r *V1) deleteTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.DeleteTaskReq

		userID, err := r.bindAuthenticatedRequest(d, &req)
		if err != nil {
			return nil, err
		}

		if err := r.tk.Delete(ctx, userID, req.ID); err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - deleteTask")

			return nil, apperror.RPC(err)
		}

		return response.DeleteStatus{Status: "deleted"}, nil
	}
}
