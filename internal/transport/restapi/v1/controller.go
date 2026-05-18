package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	n  usecase.Notification
	u  usecase.User
	tk usecase.Task
	l  logger.Interface
	v  *validator.Validate
}

func userIDFromContext(ctx *gin.Context) (string, bool) {
	userID, ok := ctx.Get("userID")
	if !ok {
		return "", false
	}

	value, ok := userID.(string)
	if !ok {
		return "", false
	}

	return value, true
}
