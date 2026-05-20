package middleware

import (
	"context"
	"strings"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const userIDKey contextKey = "userID"

const bearerParts = 2

// UserIDFromContext extracts the user ID from the context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)

	return userID, ok
}

// AuthInterceptor returns a gRPC unary interceptor for JWT authentication.
func AuthInterceptor(jwtManager *jwt.Manager) grpc.UnaryServerInterceptor {
	skipAuthMethods := map[string]bool{
		"/grpc.v1.AuthService/Register": true,
		"/grpc.v1.AuthService/Login":    true,
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if skipAuthMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, apperror.GRPCWithMessage(apperror.ErrUnauthorized, "missing metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, apperror.GRPCWithMessage(apperror.ErrUnauthorized, "missing authorization token")
		}

		parts := strings.SplitN(values[0], " ", bearerParts)
		if len(parts) != bearerParts || parts[0] != "Bearer" {
			return nil, apperror.GRPCWithMessage(apperror.ErrUnauthorized, "invalid authorization header format")
		}

		userID, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			return nil, apperror.GRPCWithMessage(apperror.ErrUnauthorized, "invalid or expired token")
		}

		ctx = context.WithValue(ctx, userIDKey, userID)

		return handler(ctx, req)
	}
}
