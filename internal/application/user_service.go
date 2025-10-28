package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend-challenge/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

var (
	generateFromPassword   = bcrypt.GenerateFromPassword
	compareHashAndPassword = bcrypt.CompareHashAndPassword
)

// UserService coordinates user use-cases.
type UserService struct {
	repo UserRepository
}

// NewUserService constructs a service with the provided repository.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterInput captures new user fields.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// UpdateInput wraps fields allowed to change.
type UpdateInput struct {
	Name  *string
	Email *string
}

// Register creates a new user with hashed password.
func (s *UserService) Register(ctx context.Context, input RegisterInput) (domain.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(strings.ToLower(input.Email))

	if err := domain.ValidateNewUser(name, email, input.Password); err != nil {
		return domain.User{}, err
	}

	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return domain.User{}, ErrDuplicateEmail
	} else if !errors.Is(err, ErrNotFound) && err != nil {
		// For errors other than not found, propagate.
		return domain.User{}, err
	}

	hashed, err := generateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	now := time.Now().UTC()
	user := domain.User{
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		CreatedAt: now,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return domain.User{}, ErrDuplicateEmail
		}
		return domain.User{}, err
	}
	return created, nil
}

// Authenticate verifies credentials and returns the user.
func (s *UserService) Authenticate(ctx context.Context, email, password string) (domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if err := domain.ValidateCredentials(email, password); err != nil {
		return domain.User{}, ErrInvalidCredentials
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return domain.User{}, ErrInvalidCredentials
		}
		return domain.User{}, err
	}

	if err := compareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domain.User{}, ErrInvalidCredentials
	}

	return user, nil
}

// Get retrieves a user by ID.
func (s *UserService) Get(ctx context.Context, id string) (domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns all users.
func (s *UserService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

// Update modifies allowed user fields.
func (s *UserService) Update(ctx context.Context, id string, input UpdateInput) (domain.User, error) {
	update := domain.UpdateUser{}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if err := domain.ValidateName(name); err != nil {
			return domain.User{}, err
		}
		update.Name = &name
	}

	if input.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*input.Email))
		if err := domain.ValidateEmail(email); err != nil {
			return domain.User{}, err
		}
		update.Email = &email
	}

	if update.Name == nil && update.Email == nil {
		return domain.User{}, ErrNoFieldsToUpdate
	}

	updated, err := s.repo.Update(ctx, id, update)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return domain.User{}, ErrDuplicateEmail
		}
		return domain.User{}, err
	}
	return updated, nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Count returns total user count.
func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
