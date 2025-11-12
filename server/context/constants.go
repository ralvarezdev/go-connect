package context

const (	
	// CtxTokenKey is the context key for the token
	CtxTokenKey = "token"
	
	// CtxTokenClaimsKey is the context key for the token claims
	CtxTokenClaimsKey = "token_claims"
	
	// CtxIssuedAccessTokenKey is the context key for the issued access token
	CtxIssuedAccessTokenKey = "issued_access_token"
	
	// CtxIssuedRefreshTokenKey is the context key for the issued refresh token
	CtxIssuedRefreshTokenKey = "issued_refresh_token"
)