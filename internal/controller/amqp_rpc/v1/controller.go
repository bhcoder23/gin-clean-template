package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	n  usecase.Notification
	u  usecase.User
	tk usecase.Task
	j  *jwt.Manager
	l  logger.Interface
	v  *validator.Validate
}
