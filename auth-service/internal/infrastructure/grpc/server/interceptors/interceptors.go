package interceptors

import (
	"context"

	"github.com/chishkin-afk/posted/auth-service/internal/domain/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(jm session.JWTManager, authRequire map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if !authRequire[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid metadata")
		}

		token := md.Get("authorization")
		if len(token) == 0 {
			return nil, status.Error(codes.InvalidArgument, "authorization is empty")
		}

		userID, err := jm.Validate(token[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, session.KeyUserID, userID)

		return handler(ctx, req)
	}
}
