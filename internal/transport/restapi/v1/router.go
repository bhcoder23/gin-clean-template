package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/middleware"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// NewRoutes -.
func NewRoutes(apiV1Group *gin.RouterGroup, n usecase.Notification, u usecase.User, tk usecase.Task, jwtManager *jwt.Manager, l logger.Interface) {
	r := &V1{n: n, u: u, tk: tk, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	r.registerAuthRoutes(apiV1Group)

	protected := apiV1Group.Group("")
	protected.Use(middleware.Auth(jwtManager))

	r.registerUserRoutes(protected)
	r.registerTaskRoutes(protected)
	r.registerNotificationRoutes(protected)
}
