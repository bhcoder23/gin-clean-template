package v1

import (
	v1 "github.com/bhcoder23/gin-clean-template/internal/transport/amqp_rpc/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

// RouterDeps groups AMQP RPC adapter dependencies.
type RouterDeps struct {
	Notification usecase.Notification
	User         usecase.User
	Task         usecase.Task
	JWTManager   *jwt.Manager
	Logger       logger.Interface
}

// NewRouter -.
func NewRouter(deps RouterDeps) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewRoutes(routes, v1.RouterDeps{
			Notification: deps.Notification,
			User:         deps.User,
			Task:         deps.Task,
			JWTManager:   deps.JWTManager,
			Logger:       deps.Logger,
		})
	}

	return routes
}
