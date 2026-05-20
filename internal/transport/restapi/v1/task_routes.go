package v1

import "github.com/gin-gonic/gin"

func (r *V1) registerTaskRoutes(protected *gin.RouterGroup) {
	taskGroup := protected.Group("/tasks")
	taskGroup.POST("", r.createTask)
	taskGroup.GET("", r.listTasks)
	taskGroup.GET("/:id", r.getTask)
	taskGroup.PUT("/:id", r.updateTask)
	taskGroup.PATCH("/:id/status", r.transitionTask)
	taskGroup.DELETE("/:id", r.deleteTask)
}
