package auth

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	goconnect "github.com/ralvarezdev/go-connect"
	goconnectserverctx "github.com/ralvarezdev/go-connect/server/context"
	gogrpc "github.com/ralvarezdev/go-grpc"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// Interceptor is the authentication interceptor
	Interceptor struct {
		validator     gojwtvalidator.Validator
		interceptions map[string]*gojwttoken.Token
		options       *Options
		logger        *slog.Logger
	}

	// Options is the options for the authentication interceptor
	Options struct {
		CustomRefreshTokenName *string
		CustomAccessTokenName  *string
		RefreshTokenFn         RefreshTokenFn
	}
)

// NewInterceptor creates a new authentication interceptor
//
// Parameters:
//
//   - validator: The JWT validator service (if nil, no validation will be done, can be used for gRPC gateways)
//   - interceptions: The map of method names to token types to intercept
//   - options: The options for the authentication interceptor (can be nil)
//   - logger: The logger for the interceptor (can be nil)
//
// Returns:
//
//   - *Interceptor: The authentication middleware
func NewInterceptor(
	validator gojwtvalidator.Validator,
	interceptions map[string]*gojwttoken.Token,
	options *Options,
	logger *slog.Logger,
) (*Interceptor, error) {
	// Check if either the validator or the gRPC interceptions is nil
	if validator == nil {
		return nil, gojwtvalidator.ErrNilValidator
	}
	if interceptions == nil {
		return nil, gogrpc.ErrNilInterceptions
	}

	// Create the logger for the interceptor
	if logger != nil {
		logger = logger.With(
			slog.String("component", "auth_interceptor"),
		)
	}

	return &Interceptor{
		validator,
		interceptions,
		options,
		logger,
	}, nil
}

// Authenticate is a unary server interceptor that extracts an authentication token
func (i Interceptor) Authenticate() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// Get the method name
			method := req.Spec().Procedure

			// Check if the method should be intercepted, if so, verify the authorization metadata is set
			interception, ok := i.interceptions[method]
			if !ok {
				// Log that the method is not intercepted
				if i.logger != nil {
					i.logger.Error(
						"Interception not found for method. Cancelling call for security reasons.",
						slog.String("method", method),
					)
				}
				return next(ctx, req)
			}
			if interception == nil {
				// Continue to the next handler if no interception is set for the method
				return next(ctx, req)
			}

			// Create a flag to indicate if the interception is for refresh tokens
			isRefreshTokenInterception := *interception == gojwttoken.RefreshToken

			// Get the cookie name from the options if set
			var (
				cookieName string
			 	customName *string
			)
			if isRefreshTokenInterception && i.options != nil {
				cookieName = goconnect.RefreshTokenCookieName
				customName = i.options.CustomRefreshTokenName
			} else if i.options != nil {
				cookieName = goconnect.AccessTokenCookieName
				customName = i.options.CustomAccessTokenName
			}

			// Try to find the token from any source
			reqHeader := req.Header()
			token, err := FindAuthorizationToken(reqHeader, customName, &cookieName)
			if err != nil {
				// Try to refresh the access token if the interception is for access tokens
				if isRefreshTokenInterception || i.options.RefreshTokenFn == nil {
					return nil, status.Error(
						codes.Unauthenticated,
						ErrUnauthenticated.Error(),
					)
				}

				// Get the refresh token from the header
				cookieName = goconnect.RefreshTokenCookieName
				refreshToken, err := FindAuthorizationToken(reqHeader, customName, &cookieName)
				if err != nil {
					return nil, status.Error(
						codes.Unauthenticated,
						ErrUnauthenticated.Error(),
					)
				}

				// Refresh the access token using the refresh token
				newRefreshToken, newAccessToken, refreshErr := i.options.RefreshTokenFn(
					ctx,
					refreshToken,
				)
				if refreshErr != nil {
					// Log the error
					if i.logger != nil {
						i.logger.Error(
							"Failed to refresh token",
							slog.String("method", method),
							slog.String("error", refreshErr.Error()),
						)
					}
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						ErrFailedToRefreshToken,
					)
				}

				// Set the new refresh and access tokens in the context as issued tokens
				ctx = goconnectserverctx.SetCtxIssuedRefreshToken(ctx, newRefreshToken)
				ctx = goconnectserverctx.SetCtxIssuedAccessToken(ctx, newAccessToken)

				// Use the new token for authentication
				if isRefreshTokenInterception {
					token = newRefreshToken
				} else {
					token = newAccessToken
				}
			}

			// Validate the token and get the validated claims
			claims, err := i.validator.ValidateClaims(ctx, token, *interception)
			if err != nil {
				return nil, status.Error(codes.Internal, gogrpc.InternalServerError)
			}

			// Set the raw token and token claims to the context
			ctx = goconnectserverctx.SetCtxToken(ctx, token)
			ctx = goconnectserverctx.SetCtxTokenClaims(ctx, claims)

			// 3. Continue to the next handler/service logic
			return next(ctx, req)
		}
	}
}
