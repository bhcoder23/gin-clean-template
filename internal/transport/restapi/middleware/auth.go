package middleware

import (
	"net/http"
	"strings"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/gin-gonic/gin"
)

const _bearerParts = 2

// Auth returns a JWT authentication middleware for Gin.
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			abortUnauthorized(ctx, "missing authorization header")

			return
		}

		parts := strings.SplitN(header, " ", _bearerParts)
		if len(parts) != _bearerParts || parts[0] != "Bearer" {
			abortUnauthorized(ctx, "invalid authorization header format")

			return
		}

		userID, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			abortUnauthorized(ctx, "invalid or expired token")

			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func abortUnauthorized(ctx *gin.Context, message string) {
	id, _ := requestid.FromContext(ctx.Request.Context())
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
		Error: response.ErrorBody{
			Code:      apperror.CodeUnauthorized,
			Message:   message,
			RequestID: id,
		},
	})
}
