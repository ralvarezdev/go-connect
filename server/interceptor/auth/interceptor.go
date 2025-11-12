package auth

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	gogrpc "github.com/ralvarezdev/go-grpc"
	gojwtgrpc "github.com/ralvarezdev/go-jwt/grpc"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	goconnect "github.com/ralvarezdev/go-connect"
	goconnectrequest "github.com/ralvarezdev/go-connect/server/request"
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
		// CustomRefreshTokenName is the custom name for the refresh token (if nil, the default name will be used)
		CustomRefreshTokenName *string

		// CustomAccessTokenName is the custom name for the access token (if nil, the default name will be used)
		CustomAccessTokenName *string

		// RefreshTokenFn is the function to refresh the access token using the refresh token
		RefreshTokenFn RefreshTokenFn
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

			// Extract the request headers from the context
			reqHeader, err := goconnectrequest.GetHeadersFromRequestContext(ctx)
			if err != nil {
				return nil, status.Error(
					codes.Unauthenticated,
					ErrUnauthenticated.Error(),
				)
			}

			// Extract the token from the request
			token, findErr := FindAuthorizationToken(reqHeader, customName, &cookieName)
			// nolint:nestif
			if findErr != nil {
				// Try to refresh the access token if the interception is for access tokens
				if isRefreshTokenInterception || i.options.RefreshTokenFn == nil {
					return nil, status.Error(
						codes.Unauthenticated,
						ErrUnauthenticated.Error(),
					)
				}

				// Get the refresh token from the header
				cookieName = goconnect.RefreshTokenCookieName
				refreshToken, findRefreshErr := FindAuthorizationToken(reqHeader, customName, &cookieName)
				if findRefreshErr != nil {
					return nil, status.Error(
						codes.Unauthenticated,
						ErrUnauthenticated.Error(),
					)
				}

				// Validate the refresh token and get the validated claims
				refreshTokenClaims, validateErr := i.validator.ValidateClaims(ctx, refreshToken, *interception)
				if validateErr != nil {
					return nil, status.Error(codes.Internal, gogrpc.InternalServerError)
				}

				// Set the raw token and token claims to the context
				ctx = gojwtgrpc.SetCtxToken(ctx, refreshToken)
				ctx = gojwtgrpc.SetCtxTokenClaims(ctx, refreshTokenClaims)

				// Refresh the access token using the refresh token claims (sho
				if refreshErr := i.options.RefreshTokenFn(
					ctx,
				); refreshErr != nil {
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
				return next(ctx, req)
			}

			// Validate the token and get the validated claims
			claims, err := i.validator.ValidateClaims(ctx, token, *interception)
			if err != nil {
				return nil, status.Error(codes.Internal, gogrpc.InternalServerError)
			}

			// Set the raw token and token claims to the context
			ctx = gojwtgrpc.SetCtxToken(ctx, token)
			ctx = gojwtgrpc.SetCtxTokenClaims(ctx, claims)

			// 3. Continue to the next handler/service logic
			return next(ctx, req)
		}
	}
}
