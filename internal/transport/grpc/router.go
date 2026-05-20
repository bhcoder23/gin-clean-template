package grpc

import (
	v1 "github.com/bhcoder23/gin-clean-template/internal/transport/grpc/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RouterDeps groups gRPC adapter dependencies.
type RouterDeps struct {
	Notification usecase.Notification
	User         usecase.User
	Task         usecase.Task
	Logger       logger.Interface
}

// NewRouter -.
func NewRouter(app *pbgrpc.Server, deps RouterDeps) {
	{
		v1.NewRoutes(app, v1.RouterDeps{
			Notification: deps.Notification,
			User:         deps.User,
			Task:         deps.Task,
			Logger:       deps.Logger,
		})
	}

	reflection.Register(app)
}
