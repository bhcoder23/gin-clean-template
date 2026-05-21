package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func extractUserID(msg *nats.Msg, jwtManager *jwt.Manager) (userID string, data json.RawMessage, err error) {
	var req request.AuthenticatedRequest

	err = json.Unmarshal(msg.Data, &req)
	if err != nil {
		return "", nil, apperror.RPC(apperror.ErrInvalidRequest)
	}

	userID, err = jwtManager.ParseToken(req.Token)
	if err != nil {
		return "", nil, apperror.RPC(apperror.ErrUnauthorized)
	}

	return userID, req.Data, nil
}

func (r *V1) bindAuthenticatedRequest(msg *nats.Msg, target any) (string, error) {
	userID, data, err := extractUserID(msg, r.j)
	if err != nil {
		return "", err
	}

	if err := r.decodeAndValidate(data, target); err != nil {
		return "", err
	}

	return userID, nil
}

func (r *V1) bindOptionalAuthenticatedRequest(msg *nats.Msg, target any) (string, error) {
	userID, data, err := extractUserID(msg, r.j)
	if err != nil {
		return "", err
	}

	if err := r.decodeOptionalAndValidate(data, target); err != nil {
		return "", err
	}

	return userID, nil
}

func authenticatedCall[Req any](
	r *V1,
	logContext string,
	call func(context.Context, string, *Req) (any, error),
) server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		var req Req

		userID, err := r.bindAuthenticatedRequest(msg, &req)
		if err != nil {
			return nil, err
		}

		result, err := call(ctx, userID, &req)
		if err != nil {
			apperror.Log(r.l, err, logContext)

			return nil, apperror.RPC(err)
		}

		return result, nil
	}
}
