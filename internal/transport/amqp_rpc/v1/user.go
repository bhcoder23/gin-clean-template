package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/internal/transport/errlog"
	"github.com/bhcoder23/gin-clean-template/internal/transport/rpcerror"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) register() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		var req request.Register

		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - register")

			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		user, err := r.u.Register(context.Background(), req.Username, req.Email, req.Password)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - register")

			return nil, rpcerror.Normalize(err)
		}

		return user, nil
	}
}

func (r *V1) login() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		var req request.Login

		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - login")

			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		token, err := r.u.Login(context.Background(), req.Email, req.Password)
		if err != nil {
			errlog.Log(r.l, err, "amqp_rpc - V1 - login")

			return nil, rpcerror.Normalize(err)
		}

		return response.Token{Token: token}, nil
	}
}
