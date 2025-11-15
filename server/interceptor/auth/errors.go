package auth

import (
	"errors"
)

const (
	ErrInterceptionNotFound = "interception not found: %s"
)

var (
	ErrUnauthenticated      = errors.New("user is unauthenticated, cannot continue with the request")
	ErrFailedToRefreshToken = errors.New("failed to refresh token")
	ErrInvalidJWTClaims     = errors.New("invalid JWT claims")
)
