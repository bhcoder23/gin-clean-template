package v1

import "github.com/gin-gonic/gin"

func (r *V1) registerPublicRoutes(public *gin.RouterGroup) {
	r.registerAuthRoutes(public)
}

func (r *V1) registerProtectedRoutes(protected *gin.RouterGroup) {
	r.registerUserRoutes(protected)
	r.registerTaskRoutes(protected)
	r.registerNotificationRoutes(protected)
}

func (r *V1) registerAuthRoutes(public *gin.RouterGroup) {
	authGroup := public.Group("/auth")
	authGroup.POST("/register", r.register)
	authGroup.POST("/login", r.login)
}

func (r *V1) registerUserRoutes(protected *gin.RouterGroup) {
	userGroup := protected.Group("/user")
	userGroup.GET("/profile", r.profile)
}

func (r *V1) registerTaskRoutes(protected *gin.RouterGroup) {
	taskGroup := protected.Group("/tasks")
	taskGroup.POST("", r.createTask)
	taskGroup.GET("", r.listTasks)
	taskGroup.GET("/:id", r.getTask)
	taskGroup.PUT("/:id", r.updateTask)
	taskGroup.PATCH("/:id/status", r.transitionTask)
	taskGroup.DELETE("/:id", r.deleteTask)
}

func (r *V1) registerNotificationRoutes(protected *gin.RouterGroup) {
	notificationGroup := protected.Group("/notifications")
	notificationGroup.GET("", r.listNotifications)
	notificationGroup.PATCH("/:id/read", r.markNotificationRead)
}
