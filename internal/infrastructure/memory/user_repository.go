package memory

import (
	"context"
	"sync"
	"time"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository is an in-memory implementation for tests.
type UserRepository struct {
	mu    sync.RWMutex
	store map[string]domain.User
}

// NewUserRepository builds an empty repository.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		store: make(map[string]domain.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range r.store {
		if v.Email == user.Email {
			return domain.User{}, application.ErrDuplicateEmail
		}
	}

	id := primitive.NewObjectID().Hex()
	user.ID = id
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
	r.store[id] = user
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.store {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, application.ErrNotFound
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.store[id]
	if !ok {
		return domain.User{}, application.ErrNotFound
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]domain.User, 0, len(r.store))
	for _, user := range r.store {
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, update domain.UpdateUser) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.store[id]
	if !ok {
		return domain.User{}, application.ErrNotFound
	}

	if update.Email != nil {
		for otherID, existing := range r.store {
			if otherID == id {
				continue
			}
			if existing.Email == *update.Email {
				return domain.User{}, application.ErrDuplicateEmail
			}
		}
		user.Email = *update.Email
	}
	if update.Name != nil {
		user.Name = *update.Name
	}

	r.store[id] = user
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return application.ErrNotFound
	}
	delete(r.store, id)
	return nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int64(len(r.store)), nil
}
