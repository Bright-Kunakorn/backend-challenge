package memory

import (
	"context"
	"errors"
	"testing"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"
)

func TestUserRepository_CreateGetListCount(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := domain.User{Name: "Test", Email: "test@example.com"}
	created, err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected ID to be set")
	}

	fetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if fetched.Email != user.Email {
		t.Fatalf("unexpected user: %+v", fetched)
	}

	fetched, err = repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}

	users, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user got %d", len(users))
	}

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count 1 got %d", count)
	}
}

func TestUserRepository_DuplicateAndUpdate(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.Create(ctx, domain.User{Name: "User", Email: "dup@example.com"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err = repo.Create(ctx, domain.User{Name: "User2", Email: "dup@example.com"}); !errors.Is(err, application.ErrDuplicateEmail) {
		t.Fatalf("expected ErrDuplicateEmail got %v", err)
	}

	user2, err := repo.Create(ctx, domain.User{Name: "User2", Email: "unique@example.com"})
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	newName := "Updated"
	updated, err := repo.Update(ctx, user2.ID, domain.UpdateUser{Name: &newName})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != newName {
		t.Fatalf("expected name %s got %s", newName, updated.Name)
	}

	newEmail := "dup@example.com"
	if _, err = repo.Update(ctx, user2.ID, domain.UpdateUser{Email: &newEmail}); !errors.Is(err, application.ErrDuplicateEmail) {
		t.Fatalf("expected duplicate email error got %v", err)
	}
}

func TestUserRepository_DeleteNotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	if err := repo.Delete(ctx, "missing"); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("expected ErrNotFound got %v", err)
	}
}
