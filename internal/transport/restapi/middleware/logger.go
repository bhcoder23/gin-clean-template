package middleware

import (
	"strconv"
	"strings"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/gin-gonic/gin"
)

func buildRequestMessage(ctx *gin.Context) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.URL.RequestURI())
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(ctx.Writer.Status()))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(ctx.Writer.Size()))

	if id, ok := requestid.FromContext(ctx.Request.Context()); ok {
		result.WriteString(" request_id=")
		result.WriteString(id)
	}

	return result.String()
}

// Logger logs request metadata after the handler chain finishes.
func Logger(l logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		l.Info("%s", buildRequestMessage(ctx))
	}
}
