package auth

import (
	"context"
)

type (
	// RefreshTokenFn defines the function signature for refreshing tokens
	RefreshTokenFn func(
		ctx context.Context,
		refreshToken string,
	) (newRefreshToken, newAccessToken string, err error)
)