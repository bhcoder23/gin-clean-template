package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
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
