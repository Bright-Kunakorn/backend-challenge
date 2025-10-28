package grpcsvc

import (
	"context"
	"net"
	"testing"
	"time"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"
	jwtinfra "backend-challenge/internal/infrastructure/jwt"
	"backend-challenge/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type repoStub struct {
	application.UserRepository
	users map[string]domain.User
}

func newRepoStub() *repoStub {
	return &repoStub{users: make(map[string]domain.User)}
}

func (r *repoStub) Create(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = "1"
	user.CreatedAt = time.Now()
	r.users[user.ID] = user
	return user, nil
}

func (r *repoStub) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return domain.User{}, application.ErrNotFound
}

func (r *repoStub) GetByID(ctx context.Context, id string) (domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return domain.User{}, application.ErrNotFound
	}
	return u, nil
}

func (r *repoStub) Count(ctx context.Context) (int64, error) {
	return int64(len(r.users)), nil
}

func TestUserServerCreateAndGet(t *testing.T) {
	repo := newRepoStub()
	service := application.NewUserService(repo)
	manager := jwtinfra.NewManager("secret", time.Hour, "issuer")

	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(AuthUnaryInterceptor(manager)))
	userServer := NewUserServer(service, manager)
	userServer.Register(server)

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("grpc server stopped: %v", err)
		}
	}()
	t.Cleanup(server.Stop)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	client := userpb.NewUserServiceClient(conn)

	createResp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{Name: "Test", Email: "test@example.com", Password: "pass12345"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if createResp.User.GetId() == "" {
		t.Fatalf("expected user id")
	}

	token, err := manager.GenerateToken(createResp.User.GetId())
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	getResp, err := client.GetUser(metadata.NewOutgoingContext(ctx, md), &userpb.GetUserRequest{Id: createResp.User.GetId()})
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if getResp.User.GetEmail() != "test@example.com" {
		t.Fatalf("unexpected user: %+v", getResp.User)
	}
}

func TestUserServerGetUnauthorized(t *testing.T) {
	repo := newRepoStub()
	service := application.NewUserService(repo)
	manager := jwtinfra.NewManager("secret", time.Hour, "issuer")

	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(AuthUnaryInterceptor(manager)))
	NewUserServer(service, manager).Register(server)
	go server.Serve(listener)
	t.Cleanup(server.Stop)

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	client := userpb.NewUserServiceClient(conn)
	_, err = client.GetUser(context.Background(), &userpb.GetUserRequest{Id: "1"})
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
}
