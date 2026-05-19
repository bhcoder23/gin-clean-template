package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) register() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.Register

		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - register")

			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		user, err := r.u.Register(ctx, req.Username, req.Email, req.Password)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - register")

			return nil, apperror.RPC(err)
		}

		return user, nil
	}
}

func (r *V1) login() server.CallHandler {
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req request.Login

		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - login")

			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		token, err := r.u.Login(ctx, req.Email, req.Password)
		if err != nil {
			apperror.Log(r.l, err, "amqp_rpc - V1 - login")

			return nil, apperror.RPC(err)
		}

		return response.Token{Token: token}, nil
	}
}
