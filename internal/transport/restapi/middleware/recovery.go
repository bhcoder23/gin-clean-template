package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/gin-gonic/gin"
)

func buildPanicMessage(ctx *gin.Context, err any) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.URL.RequestURI())
	result.WriteString(" PANIC DETECTED: ")
	fmt.Fprintf(&result, "%v\n%s\n", err, debug.Stack())

	return result.String()
}

// Recovery converts panics into logged 500 responses.
func Recovery(l logger.Interface) gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, recovered any) {
		l.Error(buildPanicMessage(ctx, recovered))
		id, _ := requestid.FromContext(ctx.Request.Context())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.NewError(
			apperror.CodeInternalServer,
			"internal server error",
			id,
		))
	})
}
