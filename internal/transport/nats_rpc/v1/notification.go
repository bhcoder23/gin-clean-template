package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func (r *V1) listNotifications() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.ListNotifications
		if len(data) > 0 {
			if err = json.Unmarshal(data, &req); err != nil {
				return nil, apperror.RPC(apperror.ErrInvalidRequest)
			}
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

		return domain.NotificationList{Notifications: notifications, Total: total}, nil
	}
}

func (r *V1) markNotificationRead() server.CallHandler {
	return func(ctx context.Context, msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.MarkNotificationRead
		if err = json.Unmarshal(data, &req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, apperror.RPC(apperror.ErrInvalidRequest)
		}

		notification, err := r.n.MarkRead(ctx, userID, req.ID)
		if err != nil {
			apperror.Log(r.l, err, "nats_rpc - V1 - markNotificationRead")

			return nil, apperror.RPC(err)
		}

		return notification, nil
	}
}
