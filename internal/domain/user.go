package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// User represents a persisted user account.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

// UserPublic is a safe projection used for API responses.
type UserPublic struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// Credentials holds login payload.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUser contains fields that may be updated.
type UpdateUser struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

var (
	// ErrInvalidName indicates the user name fails validation.
	ErrInvalidName = errors.New("name must not be empty")
	// ErrInvalidEmail indicates the email fails validation.
	ErrInvalidEmail = errors.New("email must be valid")
	// ErrInvalidPassword indicates the password fails validation.
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
)

var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

// ValidateName ensures name is not empty.
func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}
	return nil
}

// ValidateEmail ensures email is in valid format.
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(strings.TrimSpace(email)) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword enforces minimum requirements.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	return nil
}

// ValidateNewUser checks name, email, and password requirements.
func ValidateNewUser(name, email, password string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	return ValidatePassword(password)
}

// ValidateCredentials ensures email/password look reasonable.
func ValidateCredentials(email, password string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if strings.TrimSpace(password) == "" {
		return ErrInvalidPassword
	}
	return nil
}

// Sanitize converts a domain user to a public payload.
func (u User) Sanitize() UserPublic {
	return UserPublic{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
