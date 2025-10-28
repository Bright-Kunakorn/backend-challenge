package grpcsvc

import (
	"context"
	"strings"

	jwtinfra "backend-challenge/internal/infrastructure/jwt"
	"backend-challenge/internal/transport/authctx"
	"backend-challenge/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	fullMethodCreateUser = "/user.v1.UserService/CreateUser"
)

// AuthUnaryInterceptor enforces JWT authentication for gRPC calls.
func AuthUnaryInterceptor(jwtManager *jwtinfra.Manager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == fullMethodCreateUser {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization metadata")
		}

		token := values[0]
		parts := strings.SplitN(token, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
			token = parts[1]
		}

		userID, err := jwtManager.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		switch info.FullMethod {
		case "/user.v1.UserService/GetUser":
			if getReq, ok := req.(*userpb.GetUserRequest); ok && getReq.GetId() != "" && getReq.GetId() != userID {
				return nil, status.Error(codes.PermissionDenied, "forbidden")
			}
		}

		ctx = authctx.WithUserID(ctx, userID)
		return handler(ctx, req)
	}
}
