package v1

import "github.com/gin-gonic/gin"

func (r *V1) registerUserRoutes(protected *gin.RouterGroup) {
	userGroup := protected.Group("/user")
	userGroup.GET("/profile", r.profile)
}
