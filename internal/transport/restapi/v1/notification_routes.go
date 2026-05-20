package v1

import "github.com/gin-gonic/gin"

func (r *V1) registerNotificationRoutes(protected *gin.RouterGroup) {
	notificationGroup := protected.Group("/notifications")
	notificationGroup.GET("", r.listNotifications)
	notificationGroup.PATCH("/:id/read", r.markNotificationRead)
}
