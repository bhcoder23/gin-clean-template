package v1

import (
	"context"

	v1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	grpcmw "github.com/bhcoder23/gin-clean-template/internal/transport/grpc/middleware"
	"github.com/bhcoder23/gin-clean-template/internal/transport/grpc/v1/response"
)

// Register -.
func (c *AuthController) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	user, err := c.u.Register(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - Register")

		return nil, apperror.GRPC(err)
	}

	return response.NewRegisterResponse(&user), nil
}

// Login -.
func (c *AuthController) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	token, err := c.u.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - Login")

		return nil, apperror.GRPC(err)
	}

	return &v1.LoginResponse{Token: token}, nil
}

// GetProfile -.
func (c *AuthController) GetProfile(ctx context.Context, _ *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, apperror.GRPC(apperror.ErrUnauthorized)
	}

	user, err := c.u.GetUser(ctx, userID)
	if err != nil {
		apperror.Log(c.l, err, "grpc - v1 - GetProfile")

		return nil, apperror.GRPC(err)
	}

	return response.NewGetProfileResponse(&user), nil
}
