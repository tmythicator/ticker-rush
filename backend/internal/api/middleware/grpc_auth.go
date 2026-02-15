package middleware

import (
	"context"
	"strings"

	"github.com/tmythicator/ticker-rush/server/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type internalContextKey string

// UserIDContextKey is the key used to store the user ID in the context.
const UserIDContextKey = internalContextKey("userID")

// GrpcAuthInterceptor is a gRPC interceptor that validates the authorization token.
func GrpcAuthInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		accessToken := values[0]
		tokenString := strings.TrimPrefix(accessToken, "Bearer ")

		claims, err := service.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
		}

		newCtx := context.WithValue(ctx, UserIDContextKey, claims.UserID)

		return handler(newCtx, req)
	}
}
