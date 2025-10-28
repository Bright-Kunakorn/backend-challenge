package grpcsvc

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"
	jwtinfra "backend-challenge/internal/infrastructure/jwt"
	"backend-challenge/internal/transport/authctx"
	"backend-challenge/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements the gRPC UserService.
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userService *application.UserService
	jwtManager  *jwtinfra.Manager
}

// NewUserServer constructs a gRPC server wrapper.
func NewUserServer(userService *application.UserService, jwtManager *jwtinfra.Manager) *UserServer {
	return &UserServer{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// Register registers the server with a gRPC registrar.
func (s *UserServer) Register(server grpc.ServiceRegistrar) {
	userpb.RegisterUserServiceServer(server, s)
}

// CreateUser registers a new user and returns a JWT token.
func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user, err := s.userService.Register(ctx, application.RegisterInput{
		Name:     strings.TrimSpace(req.GetName()),
		Email:    strings.TrimSpace(strings.ToLower(req.GetEmail())),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	token, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &userpb.CreateUserResponse{
		User:  toProtoUser(user),
		Token: token,
	}, nil
}

// GetUser retrieves a user by ID. Requires token metadata.
func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	authID, ok := authctx.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing authentication context")
	}

	if authID != req.GetId() {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	user, err := s.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userpb.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

func toProtoUser(user domain.User) *userpb.User {
	createdAt := ""
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt.Format(time.RFC3339)
	}
	return &userpb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: createdAt,
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, application.ErrDuplicateEmail):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, application.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, application.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, application.ErrNoFieldsToUpdate),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
