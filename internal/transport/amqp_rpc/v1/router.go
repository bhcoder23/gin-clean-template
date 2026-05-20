package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/go-playground/validator/v10"
)

// NewRoutes -.
func NewRoutes(routes map[string]server.CallHandler, n usecase.Notification, u usecase.User, tk usecase.Task, j *jwt.Manager, l logger.Interface) {
	r := &V1{n: n, u: u, tk: tk, j: j, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	r.registerAuthRoutes(routes)
	r.registerNotificationRoutes(routes)
	r.registerTaskRoutes(routes)
}
