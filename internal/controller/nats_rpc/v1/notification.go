package v1

import (
	"context"
	"fmt"

	"github.com/bhcoder23/gin-clean-template/internal/controller/nats_rpc/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/entity"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func (r *V1) listNotifications() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - listNotifications - auth: %w", err)
		}

		var req request.ListNotifications
		if len(data) > 0 {
			if err = json.Unmarshal(data, &req); err != nil {
				return nil, fmt.Errorf("nats_rpc - V1 - listNotifications - json.Unmarshal: %w", err)
			}
		}

		var unreadOnly *bool
		if req.UnreadOnly {
			unreadOnly = &req.UnreadOnly
		}

		notifications, total, err := r.n.List(context.Background(), userID, unreadOnly, req.Limit, req.Offset)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - listNotifications")

			return nil, fmt.Errorf("nats_rpc - V1 - listNotifications: %w", err)
		}

		return entity.NotificationList{Notifications: notifications, Total: total}, nil
	}
}

func (r *V1) markNotificationRead() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - markNotificationRead - auth: %w", err)
		}

		var req request.MarkNotificationRead
		if err = json.Unmarshal(data, &req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - markNotificationRead - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - markNotificationRead - validation: %w", err)
		}

		notification, err := r.n.MarkRead(context.Background(), userID, req.ID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - markNotificationRead")

			return nil, fmt.Errorf("nats_rpc - V1 - markNotificationRead: %w", err)
		}

		return notification, nil
	}
}
