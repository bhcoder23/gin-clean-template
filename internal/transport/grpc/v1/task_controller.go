package v1

import (
	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// TaskController -.
type TaskController struct {
	v1.UnimplementedTaskServiceServer

	tk usecase.Task
	l  logger.Interface
	v  *validator.Validate
}
