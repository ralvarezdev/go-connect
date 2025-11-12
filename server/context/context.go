package context

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// GetCtxToken retrieves the token from the context
//
// Parameters:
//
//  - ctx: the context
//
// Returns:
//
//   - string: the token
//   - error: if there was an error retrieving the token
func GetCtxToken(ctx context.Context) (string, error) {
	// Get the token from the context
	token := ctx.Value(CtxTokenKey)
	if token == nil {
		return "", ErrNoTokenInContext
	}
	
	// Assert the token type
	tokenStr, ok := token.(string)
	if !ok {
		return "", ErrInvalidTokenInContext
	}
	return tokenStr, nil
}

// SetCtxToken sets the token in the context
// 
// Parameters:
// 
// - ctx: the context
// - token: the token
// 
// Returns:
// 
//  - context.Context: the context with the token set
func SetCtxToken(
	ctx context.Context,
	token string,
) context.Context {
	return context.WithValue(ctx, CtxTokenKey, token)
}

// GetCtxTokenClaims retrieves the token claims from the context
// 
// Parameters:
// 
// - ctx: the context
// 
// Returns:
// 
// - jwt.MapClaims: the token claims
func GetCtxTokenClaims(ctx context.Context) (jwt.MapClaims, error) {
	// Get the token claims from the context
	tokenClaims := ctx.Value(CtxTokenClaimsKey)
	if tokenClaims == nil {
		return nil, ErrNoTokenClaimsInContext
	}
	
	// Assert the token claims type
	tokenClaimsMap, ok := tokenClaims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidTokenClaimsInContext
	}
	return tokenClaimsMap, nil
}

// SetCtxTokenClaims sets the token claims in the context
// 
// Parameters:
// 
// - ctx: the context
// - tokenClaims: the token claims
// 
// Returns:
// 
// - context.Context: the context with the token claims set
func SetCtxTokenClaims(
	ctx context.Context,
	tokenClaims jwt.MapClaims,
) context.Context {
	return context.WithValue(ctx, CtxTokenClaimsKey, tokenClaims)
}

// GetCtxIssuedAccessToken retrieves the issued access token from the context
// 
// Parameters:
// 
//  - ctx: the context
// 
// Returns:
// 
//  - string: the issued access token
//  - error: if there was an error retrieving the issued access token
func GetCtxIssuedAccessToken(ctx context.Context) (string, error) {
	// Get the issued access token from the context
	issuedAccessToken := ctx.Value(CtxIssuedAccessTokenKey)
	if issuedAccessToken == nil {
		return "", ErrNoIssuedAccessTokenInContext
	}
	
	// Assert the issued access token type
	issuedAccessTokenStr, ok := issuedAccessToken.(string)
	if !ok {
		return "", ErrInvalidIssuedAccessTokenInContext
	}
	return issuedAccessTokenStr, nil
}	

// GetCtxIssuedRefreshToken retrieves the issued refresh token from the context
// 
// Parameters:
// 
// - ctx: the context
// 
// Returns:
// 
// - string: the issued refresh token
// - error: if there was an error retrieving the issued refresh token 
func GetCtxIssuedRefreshToken(ctx context.Context) (string, error) {
	// Get the issued refresh token from the context	
	issuedRefreshToken := ctx.Value(CtxIssuedRefreshTokenKey)
	if issuedRefreshToken == nil {
		return "", ErrNoIssuedRefreshTokenInContext
	}
	
	// Assert the issued refresh token type
	issuedRefreshTokenStr, ok := issuedRefreshToken.(string)
	if !ok {
		return "", ErrInvalidIssuedRefreshTokenInContext
	}
	return issuedRefreshTokenStr, nil
}

// SetCtxIssuedAccessToken sets the issued access token flag in the context
// 
// Parameters:
// 
// - ctx: the context
// - issuedAccessToken: the issued access token
// 
// Returns:
// 
// - context.Context: the context with the issued access token set
func SetCtxIssuedAccessToken(
	ctx context.Context,
	issuedAccessToken string,
) context.Context {
	return context.WithValue(ctx, CtxIssuedAccessTokenKey, issuedAccessToken)
}

// SetCtxIssuedRefreshToken sets the issued refresh token flag in the context
// 
// Parameters:
// 
// - ctx: the context
// - issuedRefreshToken: the issued refresh token
// 
// Returns:
// 
// - context.Context: the context with the issued refresh token set
func SetCtxIssuedRefreshToken(
	ctx context.Context,
	issuedRefreshToken string,
) context.Context {
	return context.WithValue(ctx, CtxIssuedRefreshTokenKey, issuedRefreshToken)
}