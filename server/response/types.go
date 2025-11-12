package headers

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"
	goconnect "github.com/ralvarezdev/go-connect"
	goconnectserverctx "github.com/ralvarezdev/go-connect/server/context"
)

type (
	// DefaultInterceptor is the default tokens response headers interceptor
	DefaultInterceptor struct {
		options *Options
	}

	// Options are the options for the tokens response headers interceptor
	Options struct {
		// CustomRefreshTokenName is the custom refresh token name to use in the response headers
		CustomRefreshTokenName *string

		// CustomAccessTokenName is the custom access token name to use in the response headers
		CustomAccessTokenName *string

		// CookieAttributes is a map of cookie attributes to use for the tokens
		CookieAttributes map[string]http.Cookie

		// CookieRefreshTokenDuration is the duration to set for the refresh token cookie
		CookieRefreshTokenDuration *time.Duration

		// CookieAccessTokenDuration is the duration to set for the access token cookie
		CookieAccessTokenDuration *time.Duration
	}
)

// NewDefaultInterceptor creates a new default tokens response headers interceptor
//
// Parameters:
//
//   - options: the options for the interceptor
//
// Returns:
//
// - *DefaultInterceptor: the default tokens response headers interceptor
func NewDefaultInterceptor(
	options *Options,
) *DefaultInterceptor {
	return &DefaultInterceptor{
		options: options,
	}
}

// InjectAccessAndRefreshToken injects the access and refresh tokens into the response headers
//
// Parameters:
//
// - ctx: the context
// - refreshToken: the refresh token
// - accessToken: the access token
//
// Returns:
//
// - error: if there was an error injecting the tokens into the response headers
func (i DefaultInterceptor) InjectAccessAndRefreshTokens(
	ctx context.Context,
	refreshToken, accessToken string,
) error {
	// Try to get the response headers from the context
	respHeader, err := GetHeadersFromRequestContext(ctx)
	if err != nil {
		return err
	}

	// Set the issued token to use as a custom header
	if i.options != nil && i.options.CustomRefreshTokenName != nil {
		respHeader.Set(*i.options.CustomRefreshTokenName, refreshToken)
	}

	// Get the cookie attributes for refresh token
	if i.options != nil && i.options.CookieAttributes != nil {
		cookieAttributes, ok := i.options.CookieAttributes[goconnect.RefreshTokenCookieName]
		if ok {
			// Set the refresh token cookie value and duration to the cookie attributes
			cookieAttributes.Value = refreshToken
			if refreshToken == "" {
				cookieAttributes.Expires = time.Time{} // Expire the cookie immediately
			}
			if i.options.CookieRefreshTokenDuration != nil {
				cookieAttributes.Expires = time.Now().Add(*i.options.CookieRefreshTokenDuration)
			}

			// Set the "Set-Cookie" header
			respHeader.Set("Set-Cookie", cookieAttributes.String())
		}
	}

	// Set the issued token to use as a custom header
	if i.options != nil && i.options.CustomAccessTokenName != nil {
		respHeader.Set(*i.options.CustomAccessTokenName, accessToken)
	}

	// Get the cookie attributes for access token
	if i.options != nil && i.options.CookieAttributes != nil {
		cookieAttributes, ok := i.options.CookieAttributes[goconnect.AccessTokenCookieName]
		if ok {
			// Set the access token cookie value and duration to the cookie attributes
			cookieAttributes.Value = accessToken
			if accessToken == "" {
				cookieAttributes.Expires = time.Time{} // Expire the cookie immediately
			}
			if i.options.CookieAccessTokenDuration != nil {
				cookieAttributes.Expires = time.Now().Add(*i.options.CookieAccessTokenDuration)
			}

			// Set the "Set-Cookie" header
			respHeader.Set("Set-Cookie", cookieAttributes.String())
		}
	}
	return nil
}

// InjectRefreshAndAccessTokensFromContext injects the refresh and access tokens from the context into the response headers
//
// Parameters:
//
//   - ctx: the context
//
// Returns:
//
// - error: if there was an error injecting the tokens into the response headers
func (i DefaultInterceptor) InjectRefreshAndAccessTokensFromContext(
	ctx context.Context,
) error {
	// Try to get the issued refresh token from the context
	refreshToken, err := goconnectserverctx.GetCtxIssuedRefreshToken(ctx)
	if err != nil {
		return err
	}

	// Try to get the issued access token from the context
	accessToken, err := goconnectserverctx.GetCtxIssuedAccessToken(ctx)
	if err != nil {
		return err
	}

	return i.InjectAccessAndRefreshTokens(
		ctx,
		refreshToken,
		accessToken,
	)
}

// CreateClientContextFromRequestContext creates a client context from the request context, injecting specified headers
// 
// Parameters:
// 
// - ctx: the request context
// - headers: the headers to inject
// 
// Returns:
// 
// - context.Context: the client context
// - error: if there was an error creating the client context
func (i DefaultInterceptor) CreateClientContextFromRequestContext(
	ctx context.Context, headers ...string,
) (context.Context, error) {
	// Get the call info from the context
	originalHeaders, err := GetHeadersFromRequestContext(ctx)
	if err != nil {
		return nil, err
	}
	
	// Create the client context
	clientCtx, callInfo := connect.NewClientContext(context.Background())
	
	// Inject the headers into the client context
	for _, header := range headers {
		// Get the original value from the header
		originalValue := originalHeaders.Get(header)
		
		// Inject the header
		callInfo.RequestHeader().Set(header, originalValue)
	}
	
	return clientCtx, err	
}
