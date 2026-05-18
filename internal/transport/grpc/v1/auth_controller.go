package v1

import (
	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// AuthController -.
type AuthController struct {
	v1.UnimplementedAuthServiceServer

	u usecase.User
	l logger.Interface
	v *validator.Validate
}
