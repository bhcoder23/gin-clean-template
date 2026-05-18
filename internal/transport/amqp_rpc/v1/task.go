package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/internal/transport/errlog"
	"github.com/bhcoder23/gin-clean-template/internal/transport/rpcerror"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) createTask() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.CreateTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		task, err := r.tk.Create(context.Background(), userID, req.Title, req.Description)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - createTask")

			return nil, rpcerror.Normalize(err)
		}

		return task, nil
	}
}

func (r *V1) getTask() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.GetTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		task, err := r.tk.Get(context.Background(), userID, req.ID)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - getTask")

			return nil, rpcerror.Normalize(err)
		}

		return task, nil
	}
}

func (r *V1) listTasks() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.ListTasks

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		var status *domain.TaskStatus

		if req.Status != "" {
			s := domain.TaskStatus(req.Status)
			status = &s
		}

		tasks, total, err := r.tk.List(context.Background(), userID, status, req.Query, req.Limit, req.Offset)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - listTasks")

			return nil, rpcerror.Normalize(err)
		}

		return response.TaskList{Tasks: tasks, Total: total}, nil
	}
}

func (r *V1) updateTask() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.UpdateTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		task, err := r.tk.Update(context.Background(), userID, req.ID, req.Title, req.Description)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - updateTask")

			return nil, rpcerror.Normalize(err)
		}

		return task, nil
	}
}

func (r *V1) transitionTask() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.TransitionTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		task, err := r.tk.Transition(context.Background(), userID, req.ID, domain.TaskStatus(req.Status))
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - transitionTask")

			return nil, rpcerror.Normalize(err)
		}

		return task, nil
	}
}

func (r *V1) deleteTask() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, err
		}

		var req request.DeleteTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		err = r.tk.Delete(context.Background(), userID, req.ID)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - deleteTask")

			return nil, rpcerror.Normalize(err)
		}

		return response.DeleteStatus{Status: "deleted"}, nil
	}
}
