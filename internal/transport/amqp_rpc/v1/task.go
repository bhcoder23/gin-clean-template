package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) createTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.CreateTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		task, err := r.tk.Create(ctx, userID, req.Title, req.Description)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - createTask")

			return nil, apperror.RPC(err)
		}

		return task, nil
	}
}

func (r *V1) getTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.GetTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		task, err := r.tk.Get(ctx, userID, req.ID)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - getTask")

			return nil, apperror.RPC(err)
		}

		return task, nil
	}
}

func (r *V1) listTasks() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.ListTasks

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
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

		return response.TaskList{Tasks: tasks, Total: total}, nil
	}
}

func (r *V1) updateTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.UpdateTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		task, err := r.tk.Update(ctx, userID, req.ID, req.Title, req.Description)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - updateTask")

			return nil, apperror.RPC(err)
		}

		return task, nil
	}
}

func (r *V1) transitionTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.TransitionTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		task, err := r.tk.Transition(ctx, userID, req.ID, domain.TaskStatus(req.Status))
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - transitionTask")

			return nil, apperror.RPC(err)
		}

		return task, nil
	}
}

func (r *V1) deleteTask() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.DeleteTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		err = r.tk.Delete(ctx, userID, req.ID)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - deleteTask")

			return nil, apperror.RPC(err)
		}

		return response.DeleteStatus{Status: "deleted"}, nil
	}
}
