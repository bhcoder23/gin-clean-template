package v1

import (
	v1 "github.com/bhcoder23/gin-clean-template/internal/controller/nats_rpc/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
)

// NewRouter -.
func NewRouter(n usecase.Notification, u usecase.User, tk usecase.Task, j *jwt.Manager, l logger.Interface) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewRoutes(routes, n, u, tk, j, l)
	}

	return routes
}
