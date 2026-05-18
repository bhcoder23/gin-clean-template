package v1

import (
	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
)

// NewNotificationRoutes -.
func NewNotificationRoutes(app *pbgrpc.Server, n usecase.Notification, l logger.Interface) {
	r := &NotificationController{n: n, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	v1.RegisterNotificationServiceServer(app, r)
}

// NewAuthRoutes -.
func NewAuthRoutes(app *pbgrpc.Server, u usecase.User, l logger.Interface) {
	r := &AuthController{u: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	v1.RegisterAuthServiceServer(app, r)
}

// NewTaskRoutes -.
func NewTaskRoutes(app *pbgrpc.Server, tk usecase.Task, l logger.Interface) {
	r := &TaskController{tk: tk, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	v1.RegisterTaskServiceServer(app, r)
}
