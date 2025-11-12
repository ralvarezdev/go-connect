package context

type (
	// ContextKey is the type for context keys
	ContextKey string
)

const (	
	// CtxTokenKey is the context key for the token
	CtxTokenKey ContextKey = "token"
	
	// CtxTokenClaimsKey is the context key for the token claims
	CtxTokenClaimsKey ContextKey = "token_claims"
	
	// CtxIssuedAccessTokenKey is the context key for the issued access token
	CtxIssuedAccessTokenKey ContextKey = "issued_access_token"
	
	// CtxIssuedRefreshTokenKey is the context key for the issued refresh token
	CtxIssuedRefreshTokenKey ContextKey = "issued_refresh_token"
)