package v1

import (
	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
)

// RouterDeps groups v1 gRPC route dependencies.
type RouterDeps struct {
	Notification usecase.Notification
	User         usecase.User
	Task         usecase.Task
	Logger       logger.Interface
}

// NewRoutes -.
func NewRoutes(app *pbgrpc.Server, deps RouterDeps) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	v1.RegisterAuthServiceServer(app, &AuthController{u: deps.User, l: deps.Logger, v: validate})
	v1.RegisterTaskServiceServer(app, &TaskController{tk: deps.Task, l: deps.Logger, v: validate})
	v1.RegisterNotificationServiceServer(app, &NotificationController{n: deps.Notification, l: deps.Logger, v: validate})
}
