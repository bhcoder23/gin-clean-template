package grpc

import (
	v1 "github.com/bhcoder23/gin-clean-template/internal/controller/grpc/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, t usecase.Translation, u usecase.User, tk usecase.Task, l logger.Interface) {
	{
		v1.NewAuthRoutes(app, u, l)
		v1.NewTaskRoutes(app, tk, l)
		v1.NewTranslationRoutes(app, t, l)
	}

	reflection.Register(app)
}
