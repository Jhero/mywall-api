package auth

import "errors"

// Authentication errors
var (
	ErrInvalidAPIKey   = errors.New("invalid API key")
	ErrInactiveAPIKey  = errors.New("API key inactive")
	ErrExpiredAPIKey   = errors.New("API key expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDatabaseError   = errors.New("database error")
	ErrUserExists      = errors.New("username or email already exists")
	ErrInvalidToken    = errors.New("invalid token")
	ErrForbidden       = errors.New("access forbidden")
	ErrUnauthorized    = errors.New("unauthorized")
)