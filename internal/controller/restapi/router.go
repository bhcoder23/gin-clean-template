package restapi

import (
	"net/http"

	"github.com/bhcoder23/gin-clean-template/config"
	_ "github.com/bhcoder23/gin-clean-template/docs" // Swagger docs.
	"github.com/bhcoder23/gin-clean-template/internal/controller/restapi/middleware"
	v1 "github.com/bhcoder23/gin-clean-template/internal/controller/restapi/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
func NewRouter(app *gin.Engine, cfg *config.Config, n usecase.Notification, u usecase.User, tk usecase.Task, jwtManager *jwt.Manager, l logger.Interface) {
	app.Use(middleware.Logger(l), middleware.Recovery(l))

	if cfg.Metrics.Enabled {
		app.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	if cfg.Swagger.Enabled {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	app.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	apiV1Group := app.Group("/v1")
	v1.NewRoutes(apiV1Group, n, u, tk, jwtManager, l)
}
