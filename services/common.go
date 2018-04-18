package services

import "errors"

var (
	// ErrNotAllowed occurs when user does not have sufficient auth level
	ErrNotAllowed = errors.New("user is not allowed")

	// ErrForbidden occurs when a user is forbidden from doing something
	ErrForbidden = errors.New("user is forbidden")
)
