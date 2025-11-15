package goconnect

import (
	"errors"
)

var (
	ErrInvalidAuthorization = errors.New("invalid authorization HTTP header or gRPC metadata")
	ErrMissingAuthorization = errors.New("missing authorization HTTP header or gRPC metadata")
	ErrClientIPNotFound = errors.New("client IP address not found in request context")
	ErrInternalServerError = errors.New("internal server error")
)
