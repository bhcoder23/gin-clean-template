package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/bhcoder23/gin-clean-template/config"
	_ "github.com/bhcoder23/gin-clean-template/docs" // Swagger docs.
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/middleware"
	v1 "github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const _defaultReadinessTimeout = 2 * time.Second

type readinessResponse struct {
	Error string `json:"error"`
}

// RouterDeps groups REST adapter dependencies.
type RouterDeps struct {
	Notification usecase.Notification
	User         usecase.User
	Task         usecase.Task
	JWTManager   *jwt.Manager
	Logger       logger.Interface
}

type routerOptions struct {
	readinessCheck   func(context.Context) error
	readinessTimeout time.Duration
}

// Option configures REST router infrastructure hooks.
type Option func(*routerOptions)

// ReadinessCheck registers a dependency readiness check.
func ReadinessCheck(check func(context.Context) error) Option {
	return func(o *routerOptions) {
		o.readinessCheck = check
	}
}

// NewRouter -.
// Swagger spec:
//
//	@title       Gin Clean Template API
//	@description Multi-domain clean architecture template with notifications, user, and task management
//	@version     1.0
//	@host        localhost:8080
//	@BasePath    /v1
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
func NewRouter(app *gin.Engine, cfg *config.Config, deps RouterDeps, opts ...Option) {
	options := routerOptions{readinessTimeout: _defaultReadinessTimeout}
	for _, opt := range opts {
		opt(&options)
	}

	app.Use(middleware.RequestID())

	if cfg.Trace.Enabled {
		app.Use(middleware.Tracing(cfg.Trace.ServiceName))
	}

	app.Use(middleware.Logger(deps.Logger), middleware.Recovery(deps.Logger))

	if cfg.Metrics.Enabled {
		app.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	if cfg.Swagger.Enabled {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	app.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	app.GET("/readyz", func(ctx *gin.Context) {
		if options.readinessCheck == nil {
			ctx.Status(http.StatusOK)

			return
		}

		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), options.readinessTimeout)
		defer cancel()

		if err := options.readinessCheck(checkCtx); err != nil {
			deps.Logger.Error(err, "restapi - readyz")
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, readinessResponse{Error: "service unavailable"})

			return
		}

		ctx.Status(http.StatusOK)
	})

	apiV1Group := app.Group("/v1")
	v1.NewRoutes(apiV1Group, v1.RouterDeps{
		Notification: deps.Notification,
		User:         deps.User,
		Task:         deps.Task,
		JWTManager:   deps.JWTManager,
		Logger:       deps.Logger,
	})
}
