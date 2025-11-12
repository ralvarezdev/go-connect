package context

import (
	"context"
)

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