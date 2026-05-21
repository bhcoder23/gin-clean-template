package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func extractUserID(d *amqp.Delivery, jwtManager *jwt.Manager) (userID string, data json.RawMessage, err error) {
	var req request.AuthenticatedRequest

	err = json.Unmarshal(d.Body, &req)
	if err != nil {
		return "", nil, apperror.RPC(apperror.ErrInvalidRequest)
	}

	userID, err = jwtManager.ParseToken(req.Token)
	if err != nil {
		return "", nil, apperror.RPC(apperror.ErrUnauthorized)
	}

	return userID, req.Data, nil
}

func (r *V1) bindAuthenticatedRequest(d *amqp.Delivery, target any) (string, error) {
	userID, data, err := extractUserID(d, r.j)
	if err != nil {
		return "", err
	}

	if err := r.decodeAndValidate(data, target); err != nil {
		return "", err
	}

	return userID, nil
}

func (r *V1) bindOptionalAuthenticatedRequest(d *amqp.Delivery, target any) (string, error) {
	userID, data, err := extractUserID(d, r.j)
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
	return func(ctx context.Context, d *amqp.Delivery) (any, error) {
		var req Req

		userID, err := r.bindAuthenticatedRequest(d, &req)
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
