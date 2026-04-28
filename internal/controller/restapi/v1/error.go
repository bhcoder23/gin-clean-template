package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/controller/restapi/v1/response"
	"github.com/gin-gonic/gin"
)

func errorResponse(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(code, response.Error{Error: msg})
}
