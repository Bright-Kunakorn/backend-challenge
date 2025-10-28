package userpb

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type fakeUserService struct{}

func (f *fakeUserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{
		User:  &User{Id: "1", Name: req.GetName(), Email: req.GetEmail()},
		Token: "token",
	}, nil
}

func (f *fakeUserService) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	return &GetUserResponse{User: &User{Id: req.GetId(), Name: "Test", Email: "test@example.com"}}, nil
}

func (f *fakeUserService) mustEmbedUnimplementedUserServiceServer() {}

func TestUserServiceClientServer(t *testing.T) {
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	RegisterUserServiceServer(server, &fakeUserService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("server exited: %v", err)
		}
	}()
	t.Cleanup(server.Stop)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	client := NewUserServiceClient(conn)

	createResp, err := client.CreateUser(ctx, &CreateUserRequest{Name: "Test", Email: "test@example.com", Password: "pass"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if createResp.User.GetId() != "1" || createResp.Token == "" {
		t.Fatalf("unexpected create response: %+v", createResp)
	}

	getResp, err := client.GetUser(ctx, &GetUserRequest{Id: "1"})
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if getResp.User.GetEmail() != "test@example.com" {
		t.Fatalf("unexpected get response: %+v", getResp)
	}
}
