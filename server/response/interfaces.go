package response

import (
	"context"

	"connectrpc.com/connect"
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
		InjectHeadersFromCallInfo(
			ctx context.Context,
			callInfo connect.CallInfo,
		) error
	}
)
