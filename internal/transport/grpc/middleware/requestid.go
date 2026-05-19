package middleware

import (
	"context"

	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIDInterceptor carries request ids through gRPC metadata and context.
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		id := ""

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get(requestid.MetadataKey)
			if len(values) > 0 {
				id = values[0]
			}
		}

		id = requestid.Normalize(id)
		ctx = requestid.WithContext(ctx, id)

		if err := grpc.SetHeader(ctx, metadata.Pairs(requestid.MetadataKey, id)); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
