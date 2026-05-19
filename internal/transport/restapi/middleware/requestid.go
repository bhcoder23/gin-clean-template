package middleware

import (
	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/gin-gonic/gin"
)

// RequestID ensures each HTTP request has a correlation id.
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := requestid.Normalize(ctx.GetHeader(requestid.Header))
		ctx.Request = ctx.Request.WithContext(requestid.WithContext(ctx.Request.Context(), id))
		ctx.Header(requestid.Header, id)
		ctx.Set(requestid.MetadataKey, id)
		ctx.Next()
	}
}
