package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/middleware"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RouterDeps groups v1 REST route dependencies.
type RouterDeps struct {
	Notification usecase.Notification
	User         usecase.User
	Task         usecase.Task
	JWTManager   *jwt.Manager
	Logger       logger.Interface
}

// NewRoutes -.
func NewRoutes(apiV1Group *gin.RouterGroup, deps RouterDeps) {
	r := &V1{
		n:  deps.Notification,
		u:  deps.User,
		tk: deps.Task,
		l:  deps.Logger,
		v:  validator.New(validator.WithRequiredStructEnabled()),
	}

	public := apiV1Group.Group("")
	protected := apiV1Group.Group("")
	protected.Use(middleware.Auth(deps.JWTManager))

	r.registerPublicRoutes(public)
	r.registerProtectedRoutes(protected)
}
