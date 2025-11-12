package context

import (
	"errors"
)

var (
	ErrNoTokenInContext                   = errors.New("no token in context")
	ErrInvalidTokenInContext              = errors.New("invalid token in context")
	ErrNoTokenClaimsInContext             = errors.New("no token claims in context")
	ErrInvalidTokenClaimsInContext        = errors.New("invalid token claims in context")
	ErrNoIssuedAccessTokenInContext       = errors.New("no issued access token in context")
	ErrInvalidIssuedAccessTokenInContext  = errors.New("invalid issued access token in context")
	ErrNoIssuedRefreshTokenInContext      = errors.New("no issued refresh token in context")
	ErrInvalidIssuedRefreshTokenInContext = errors.New("invalid issued refresh token in context")
)
