package domain

import (
	"testing"
	"time"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "John Doe"},
		{name: "trimmed valid", input: "  Jane  "},
		{name: "empty", input: "   ", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateName(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error got %v", err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "john@example.com"},
		{name: "mixed case", input: "Jane.Doe+test@Example.co"},
		{name: "missing at", input: "janeexample.com", wantErr: true},
		{name: "missing domain", input: "jane@", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error got %v", err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "long-enough"},
		{name: "short", input: "short", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error got %v", err)
			}
		})
	}
}

func TestValidateNewUser(t *testing.T) {
	if err := ValidateNewUser("Name", "user@example.com", "password123"); err != nil {
		t.Fatalf("expected nil error got %v", err)
	}

	err := ValidateNewUser("", "user@example.com", "password123")
	if err != ErrInvalidName {
		t.Fatalf("expected ErrInvalidName got %v", err)
	}

	err = ValidateNewUser("Name", "invalid", "password123")
	if err != ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail got %v", err)
	}

	err = ValidateNewUser("Name", "user@example.com", "short")
	if err != ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword got %v", err)
	}
}

func TestValidateCredentials(t *testing.T) {
	if err := ValidateCredentials("user@example.com", "password123"); err != nil {
		t.Fatalf("expected nil error got %v", err)
	}

	if err := ValidateCredentials("bad", "password123"); err != ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail got %v", err)
	}

	if err := ValidateCredentials("user@example.com", "   "); err != ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword got %v", err)
	}
}

func TestSanitize(t *testing.T) {
	created := time.Now().UTC()
	user := User{
		ID:        "abc123",
		Name:      "John",
		Email:     "john@example.com",
		Password:  "hashed",
		CreatedAt: created,
	}

	public := user.Sanitize()

	if public.ID != user.ID {
		t.Fatalf("expected ID %s got %s", user.ID, public.ID)
	}
	if public.Name != user.Name {
		t.Fatalf("expected Name %s got %s", user.Name, public.Name)
	}
	if public.Email != user.Email {
		t.Fatalf("expected Email %s got %s", user.Email, public.Email)
	}
	if !public.CreatedAt.Equal(created) {
		t.Fatalf("expected CreatedAt %v got %v", created, public.CreatedAt)
	}
}
