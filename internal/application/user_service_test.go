package application_test

import (
	"context"
	"testing"

	"backend-challenge/internal/application"
	"backend-challenge/internal/infrastructure/memory"

	"github.com/stretchr/testify/require"
)

func newService() (*application.UserService, context.Context) {
	repo := memory.NewUserRepository()
	svc := application.NewUserService(repo)
	return svc, context.Background()
}

func TestRegisterAndAuthenticate(t *testing.T) {
	service, ctx := newService()

	user, err := service.Register(ctx, application.RegisterInput{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "supersecret",
	})
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.NotEmpty(t, user.Password)
	require.NotEqual(t, "supersecret", user.Password, "password should be hashed")

	authenticated, err := service.Authenticate(ctx, "jane@example.com", "supersecret")
	require.NoError(t, err)
	require.Equal(t, user.ID, authenticated.ID)
}

func TestRegisterDuplicateEmail(t *testing.T) {
	service, ctx := newService()

	_, err := service.Register(ctx, application.RegisterInput{
		Name:     "User A",
		Email:    "duplicate@example.com",
		Password: "supersecret",
	})
	require.NoError(t, err)

	_, err = service.Register(ctx, application.RegisterInput{
		Name:     "User B",
		Email:    "duplicate@example.com",
		Password: "anothersecret",
	})
	require.ErrorIs(t, err, application.ErrDuplicateEmail)
}

func TestAuthenticateInvalidPassword(t *testing.T) {
	service, ctx := newService()

	_, err := service.Register(ctx, application.RegisterInput{
		Name:     "John Smith",
		Email:    "john@example.com",
		Password: "goodpassword",
	})
	require.NoError(t, err)

	_, err = service.Authenticate(ctx, "john@example.com", "wrongpassword")
	require.ErrorIs(t, err, application.ErrInvalidCredentials)
}

func TestUpdateAndDelete(t *testing.T) {
	service, ctx := newService()

	user, err := service.Register(ctx, application.RegisterInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	newName := "Alice Smith"
	updated, err := service.Update(ctx, user.ID, application.UpdateInput{
		Name: &newName,
	})
	require.NoError(t, err)
	require.Equal(t, newName, updated.Name)

	err = service.Delete(ctx, user.ID)
	require.NoError(t, err)

	_, err = service.Get(ctx, user.ID)
	require.ErrorIs(t, err, application.ErrNotFound)
}

func TestUpdateDuplicateEmail(t *testing.T) {
	service, ctx := newService()

	u1, err := service.Register(ctx, application.RegisterInput{
		Name:     "Alpha",
		Email:    "alpha@example.com",
		Password: "password123",
	})
	require.NoError(t, err)
	_, err = service.Register(ctx, application.RegisterInput{
		Name:     "Beta",
		Email:    "beta@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	_, err = service.Update(ctx, u1.ID, application.UpdateInput{
		Email: strPtr("beta@example.com"),
	})
	require.ErrorIs(t, err, application.ErrDuplicateEmail)
}

func strPtr(value string) *string {
	return &value
}
