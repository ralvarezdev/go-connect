package auth

import (
	"context"
)

type (
	// RefreshTokenFn defines the function signature for refreshing tokens.
	// NOTE: It should obtain the refresh token from the context and it should set the issued tokens in the context.
	RefreshTokenFn func(
		ctx context.Context,
	) error
)
