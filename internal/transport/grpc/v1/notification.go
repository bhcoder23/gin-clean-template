package v1

import (
	"context"

	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	grpcmw "github.com/bhcoder23/gin-clean-template/internal/transport/grpc/middleware"
	"github.com/bhcoder23/gin-clean-template/internal/transport/grpc/v1/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListNotifications -.
func (c *NotificationController) ListNotifications(ctx context.Context, req *v1.ListNotificationsRequest) (*v1.ListNotificationsResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	var unreadOnly *bool

	if req.GetUnreadOnly() {
		value := true
		unreadOnly = &value
	}

	notifications, total, err := c.n.List(ctx, userID, unreadOnly, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - ListNotifications")

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return response.NewListNotificationsResponse(notifications, total), nil
}

// MarkNotificationRead -.
func (c *NotificationController) MarkNotificationRead(ctx context.Context, req *v1.MarkNotificationReadRequest) (*v1.NotificationResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	notification, err := c.n.MarkRead(ctx, userID, req.GetId())
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - MarkNotificationRead")

		return nil, apperror.GRPC(err)
	}

	return response.NewNotificationResponse(&notification), nil
}
