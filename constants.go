package goconnect

const (
	// AuthorizationKey is the key for the authorization header and gRPC metadata
	AuthorizationKey = "Authorization"

	// RefreshTokenCookieName is the name of the refresh token cookie
	RefreshTokenCookieName = "X-Refresh-Token"

	// AccessTokenCookieName is the name of the access token cookie
	AccessTokenCookieName = "X-Access-Token"
)
