package services

import "errors"

var (
	// ErrNotAllowed occurs when user does not have sufficient auth level
	ErrNotAllowed = errors.New("user is not allowed")
)
