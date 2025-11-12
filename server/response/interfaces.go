package headers

import (
	"context"
)

type (
	// Injector is used to inject headers into the response headers
	Injector interface {
		InjectAccessAndRefreshTokens(
			ctx context.Context,
			refreshToken,
			accessToken string,
		) error
		InjectAccessAndRefreshTokensFromContext(
			ctx context.Context,
		) error
		CreateClientContextFromRequestContext(
			ctx context.Context, headers ...string,
		) (context.Context, error)
	}
)
