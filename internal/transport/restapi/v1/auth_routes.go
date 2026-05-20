package v1

import "github.com/gin-gonic/gin"

func (r *V1) registerAuthRoutes(apiV1Group *gin.RouterGroup) {
	authGroup := apiV1Group.Group("/auth")
	authGroup.POST("/register", r.register)
	authGroup.POST("/login", r.login)
}
