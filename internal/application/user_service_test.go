package application_test

import (
	"context"
	"errors"
	"testing"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"
	"backend-challenge/internal/infrastructure/memory"

	"github.com/stretchr/testify/require"
)

func newService() (*application.UserService, context.Context) {
	repo := memory.NewUserRepository()
	return application.NewUserService(repo), context.Background()
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
	require.NotEqual(t, "supersecret", user.Password)

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

func TestRegisterValidationErrors(t *testing.T) {
	service, ctx := newService()

	_, err := service.Register(ctx, application.RegisterInput{Name: " ", Email: "user@example.com", Password: "password123"})
	require.ErrorIs(t, err, domain.ErrInvalidName)

	_, err = service.Register(ctx, application.RegisterInput{Name: "Name", Email: "bad-email", Password: "password123"})
	require.ErrorIs(t, err, domain.ErrInvalidEmail)

	_, err = service.Register(ctx, application.RegisterInput{Name: "Name", Email: "user@example.com", Password: "short"})
	require.ErrorIs(t, err, domain.ErrInvalidPassword)
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
	updated, err := service.Update(ctx, user.ID, application.UpdateInput{Name: &newName})
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

	_, err = service.Update(ctx, u1.ID, application.UpdateInput{Email: strPtr("beta@example.com")})
	require.ErrorIs(t, err, application.ErrDuplicateEmail)
}

func TestRegisterRepositoryErrors(t *testing.T) {
	ctx := context.Background()
	unexpectedErr := errors.New("unexpected")
	repo := &stubRepo{
		getByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{}, unexpectedErr
		},
	}
	service := application.NewUserService(repo)

	_, err := service.Register(ctx, application.RegisterInput{Name: "Name", Email: "user@example.com", Password: "password123"})
	require.ErrorIs(t, err, unexpectedErr)

	repo.getByEmail = func(context.Context, string) (domain.User, error) {
		return domain.User{}, application.ErrNotFound
	}
	repo.createFn = func(context.Context, domain.User) (domain.User, error) {
		return domain.User{}, application.ErrDuplicateEmail
	}

	_, err = service.Register(ctx, application.RegisterInput{Name: "Name", Email: "user@example.com", Password: "password123"})
	require.ErrorIs(t, err, application.ErrDuplicateEmail)

	repo.createFn = func(context.Context, domain.User) (domain.User, error) {
		return domain.User{}, errors.New("create failed")
	}

	_, err = service.Register(ctx, application.RegisterInput{Name: "Name", Email: "user2@example.com", Password: "password123"})
	require.ErrorContains(t, err, "create failed")
}

func TestRegisterHashFailure(t *testing.T) {
	ctx := context.Background()
	restore := application.OverrideHashFuncForTests(func([]byte, int) ([]byte, error) {
		return nil, errors.New("hash failed")
	})
	defer restore()

	repo := &stubRepo{
		getByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{}, application.ErrNotFound
		},
	}
	service := application.NewUserService(repo)

	_, err := service.Register(ctx, application.RegisterInput{Name: "Name", Email: "user@example.com", Password: "password123"})
	require.Error(t, err)
}

func TestAuthenticateRepositoryErrors(t *testing.T) {
	ctx := context.Background()
	unexpectedErr := errors.New("repo error")
	repo := &stubRepo{
		getByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{}, unexpectedErr
		},
	}
	service := application.NewUserService(repo)

	_, err := service.Authenticate(ctx, "user@example.com", "password123")
	require.ErrorIs(t, err, unexpectedErr)
}

func TestAuthenticateInvalidInputsAndNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &stubRepo{}
	service := application.NewUserService(repo)

	_, err := service.Authenticate(ctx, "bad-email", "password123")
	require.ErrorIs(t, err, application.ErrInvalidCredentials)

	_, err = service.Authenticate(ctx, "user@example.com", "password123")
	require.ErrorIs(t, err, application.ErrInvalidCredentials)
}

func TestListAndCount(t *testing.T) {
	ctx := context.Background()
	repo := &stubRepo{
		listFn: func(context.Context) ([]domain.User, error) {
			return []domain.User{{ID: "1"}}, nil
		},
		countFn: func(context.Context) (int64, error) {
			return 42, nil
		},
	}
	service := application.NewUserService(repo)

	users, err := service.List(ctx)
	require.NoError(t, err)
	require.Len(t, users, 1)

	count, err := service.Count(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(42), count)
}

func TestUpdateNoFields(t *testing.T) {
	ctx := context.Background()
	repo := &stubRepo{}
	service := application.NewUserService(repo)

	_, err := service.Update(ctx, "id", application.UpdateInput{})
	require.ErrorIs(t, err, application.ErrNoFieldsToUpdate)
}

func TestDeleteError(t *testing.T) {
	ctx := context.Background()
	errDelete := errors.New("delete failed")
	repo := &stubRepo{
		deleteFn: func(context.Context, string) error {
			return errDelete
		},
	}
	service := application.NewUserService(repo)

	err := service.Delete(ctx, "id")
	require.ErrorIs(t, err, errDelete)
}

func TestUpdateValidationAndRepositoryErrors(t *testing.T) {
	ctx := context.Background()
	repo := &stubRepo{}
	service := application.NewUserService(repo)

	_, err := service.Update(ctx, "id", application.UpdateInput{Name: strPtr("   ")})
	require.ErrorIs(t, err, domain.ErrInvalidName)

	_, err = service.Update(ctx, "id", application.UpdateInput{Email: strPtr("invalid")})
	require.ErrorIs(t, err, domain.ErrInvalidEmail)

	repo.updateFn = func(context.Context, string, domain.UpdateUser) (domain.User, error) {
		return domain.User{}, errors.New("update failed")
	}

	name := "New"
	_, err = service.Update(ctx, "id", application.UpdateInput{Name: &name})
	require.ErrorContains(t, err, "update failed")
}

func strPtr(value string) *string {
	return &value
}

type stubRepo struct {
	createFn   func(context.Context, domain.User) (domain.User, error)
	getByEmail func(context.Context, string) (domain.User, error)
	getByID    func(context.Context, string) (domain.User, error)
	listFn     func(context.Context) ([]domain.User, error)
	updateFn   func(context.Context, string, domain.UpdateUser) (domain.User, error)
	deleteFn   func(context.Context, string) error
	countFn    func(context.Context) (int64, error)
}

func (s *stubRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	if s.createFn != nil {
		return s.createFn(ctx, user)
	}
	return domain.User{}, nil
}

func (s *stubRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if s.getByEmail != nil {
		return s.getByEmail(ctx, email)
	}
	return domain.User{}, application.ErrNotFound
}

func (s *stubRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	if s.getByID != nil {
		return s.getByID(ctx, id)
	}
	return domain.User{}, application.ErrNotFound
}

func (s *stubRepo) List(ctx context.Context) ([]domain.User, error) {
	if s.listFn != nil {
		return s.listFn(ctx)
	}
	return nil, nil
}

func (s *stubRepo) Update(ctx context.Context, id string, update domain.UpdateUser) (domain.User, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, update)
	}
	return domain.User{}, nil
}

func (s *stubRepo) Delete(ctx context.Context, id string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

func (s *stubRepo) Count(ctx context.Context) (int64, error) {
	if s.countFn != nil {
		return s.countFn(ctx)
	}
	return 0, nil
}
