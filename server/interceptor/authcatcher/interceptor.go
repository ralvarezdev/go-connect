package auth

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"

	goconnect "github.com/ralvarezdev/go-connect"
	goconnectserverctx "github.com/ralvarezdev/go-connect/server/context"
)

type (
	// Interceptor is the authentication catcher interceptor
	Interceptor struct{
		options *Options
	}
	
	// Options are the options for the authentication catcher interceptor
	Options struct {
		CustomRefreshTokenName *string
		CustomAccessTokenName  *string
		CookieRefreshTokenDuration  *time.Duration
		CookieAccessTokenDuration   *time.Duration
		CookieAttributes map[string]http.Cookie
	}
)

// NewInterceptor creates a new authentication catcher interceptor
// 
// Parameters:
// 
//  - options: the options for the interceptor
//
// Returns:
// 
// - *Interceptor: the authentication catcher interceptor
func NewInterceptor(options *Options) *Interceptor {
	return &Interceptor{
		options: options,
	}
}

// Catch is the unary interceptor that catches the issued tokens from the context and sets them in the response metadata and cookies
func (i Interceptor) Catch() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// Call the next handler in the chain
			resp, err := next(ctx, req)
			if err != nil {
				// If the login failed, do nothing and return the error.
				return resp, err 
			}
			
			// Get the response header
			respHeader := resp.Header()

			// Try to get the issued refresh and access tokens from the context
			refreshToken, err := goconnectserverctx.GetCtxIssuedRefreshToken(ctx)
			if err == nil {
				// Set the issued token to use as cookie and header
				if i.options != nil && i.options.CustomRefreshTokenName != nil {
					respHeader.Set(*i.options.CustomRefreshTokenName, refreshToken)
				}
				
				// Get the cookie attributes for refresh token
				if i.options != nil && i.options.CookieAttributes != nil {
					cookieAttributes, ok := i.options.CookieAttributes[goconnect.RefreshTokenCookieName]
					if ok {
						// Set the refresh token cookie value and duration to the cookie attributes
						cookieAttributes.Value = refreshToken
						if i.options.CookieRefreshTokenDuration != nil {
							cookieAttributes.Expires = time.Now().Add(*i.options.CookieRefreshTokenDuration)
						}
						
						// Set the "Set-Cookie" header
						respHeader.Set("Set-Cookie", cookieAttributes.String())
					}
				}
			}
				
			// Try to get the issued access token from the context
			accessToken, err := goconnectserverctx.GetCtxIssuedAccessToken(ctx)
			if err == nil {
				// Set the issued token to use as cookie and header
				if i.options != nil && i.options.CustomAccessTokenName != nil {
					respHeader.Set(*i.options.CustomAccessTokenName, accessToken)
				}
				
				// Get the cookie attributes for access token
				if i.options != nil && i.options.CookieAttributes != nil {
					cookieAttributes, ok := i.options.CookieAttributes[goconnect.AccessTokenCookieName]
					if ok {
						// Set the access token cookie value and duration to the cookie attributes
						cookieAttributes.Value = accessToken
						if i.options.CookieAccessTokenDuration != nil {
							cookieAttributes.Expires = time.Now().Add(*i.options.CookieAccessTokenDuration)
						}
						
						// Set the "Set-Cookie" header
						respHeader.Set("Set-Cookie", cookieAttributes.String())
					}
				}
			}
		
			return resp, err
		}
	}
}