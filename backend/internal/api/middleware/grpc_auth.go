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

const UserIDContextKey = internalContextKey("userID")

func GrpcAuthInterceptor(
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

	claims, err := service.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	newCtx := context.WithValue(ctx, UserIDContextKey, claims.UserID)

	return handler(newCtx, req)
}
