package v1

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/errlog"
	"github.com/bhcoder23/gin-clean-template/internal/transport/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/rpcerror"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func (r *V1) listNotifications() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.ListNotifications
		if len(data) > 0 {
			if err = json.Unmarshal(data, &req); err != nil {
				return nil, rpcerror.ErrInvalidRequest
			}
		}

		var unreadOnly *bool
		if req.UnreadOnly {
			unreadOnly = &req.UnreadOnly
		}

		notifications, total, err := r.n.List(context.Background(), userID, unreadOnly, req.Limit, req.Offset)
		if err != nil {
			errlog.Log(r.l, err, "nats_rpc - V1 - listNotifications")

			return nil, rpcerror.Normalize(err)
		}

		return domain.NotificationList{Notifications: notifications, Total: total}, nil
	}
}

func (r *V1) markNotificationRead() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, err
		}

		var req request.MarkNotificationRead
		if err = json.Unmarshal(data, &req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		if err = r.v.Struct(req); err != nil {
			return nil, rpcerror.ErrInvalidRequest
		}

		notification, err := r.n.MarkRead(context.Background(), userID, req.ID)
		if err != nil {
			errlog.Log(r.l, err, "nats_rpc - V1 - markNotificationRead")

			return nil, rpcerror.Normalize(err)
		}

		return notification, nil
	}
}
