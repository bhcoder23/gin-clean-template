package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func (r *V1) register() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		var req request.RegisterReq

		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - register")

			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		user, err := r.u.Register(ctx, req.Username, req.Email, req.Password)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - register")

			return nil, apperror.RPC(err)
		}

		return response.NewUserResp(user), nil
	}
}

func (r *V1) login() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		var req request.LoginReq

		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - login")

			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		token, err := r.u.Login(ctx, req.Email, req.Password)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - login")

			return nil, apperror.RPC(err)
		}

		return response.TokenResp{Token: token}, nil
	}
}
