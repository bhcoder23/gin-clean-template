package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/nats-io/nats.go"
)

func (r *V1) listNotifications() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		var req request.ListNotificationsReq

		userID, err := r.bindOptionalAuthenticatedRequest(msg, &req)
		if err != nil {
			return nil, err
		}

		var unreadOnly *bool
		if req.UnreadOnly {
			unreadOnly = &req.UnreadOnly
		}

		notifications, total, err := r.n.List(ctx, userID, unreadOnly, req.Limit, req.Offset)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - listNotifications")

			return nil, apperror.RPC(err)
		}

		return response.NewListNotificationsResp(notifications, total), nil
	}
}

//nolint:dupl // RPC handlers stay explicit per route; shared binding/error mapping is already factored.
func (r *V1) markNotificationRead() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		var req request.MarkNotificationReadReq

		userID, err := r.bindAuthenticatedRequest(msg, &req)
		if err != nil {
			return nil, err
		}

		notification, err := r.n.MarkRead(ctx, userID, req.ID)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - markNotificationRead")

			return nil, apperror.RPC(err)
		}

		return response.NewNotificationResp(&notification), nil
	}
}
