package application

import "errors"

var (
	// ErrNotFound indicates the requested entity does not exist.
	ErrNotFound = errors.New("user not found")
	// ErrDuplicateEmail indicates the email already exists.
	ErrDuplicateEmail = errors.New("email already in use")
	// ErrInvalidCredentials indicates login failed.
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrNoFieldsToUpdate indicates update payload missing fields.
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)
