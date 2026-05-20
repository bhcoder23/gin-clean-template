package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/go-playground/validator/v10"
)

// RouterDeps groups v1 NATS RPC route dependencies.
type RouterDeps struct {
	Notification usecase.Notification
	User         usecase.User
	Task         usecase.Task
	JWTManager   *jwt.Manager
	Logger       logger.Interface
}

// NewRoutes -.
func NewRoutes(routes map[string]server.CallHandler, deps RouterDeps) {
	r := &V1{
		n:  deps.Notification,
		u:  deps.User,
		tk: deps.Task,
		j:  deps.JWTManager,
		l:  deps.Logger,
		v:  validator.New(validator.WithRequiredStructEnabled()),
	}

	r.registerAuthRoutes(routes)
	r.registerNotificationRoutes(routes)
	r.registerTaskRoutes(routes)
}
