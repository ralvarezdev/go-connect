package goconnect

const (
	// AuthorizationKey is the key for the authorization header and gRPC metadata
	AuthorizationKey = "Authorization"

	// RefreshTokenCookieName is the name of the refresh token cookie
	RefreshTokenCookieName = "X-Refresh-Token"

	// AccessTokenCookieName is the name of the access token cookie
	AccessTokenCookieName = "X-Access-Token"
	
	// AccessTokenKey is the key for the access token
	AccessTokenKey = "X-Access-Token"
	
	// RefreshTokenKey is the key for the refresh token
	RefreshTokenKey = "X-Refresh-Token"
	
	// XForwardedForKey is the key for the X-Forwarded-For header
	XForwardedForKey = "X-Forwarded-For"

	// RemoteAddrKey is the key for the remote address
	RemoteAddrKey = "Remote-Addr"
	
	// XRealIPKey is the key for the X-Real-IP header
	XRealIPKey = "X-Real-IP"
)

var (
	// AuthHeaders are the headers used for authentication
	AuthHeaders = []string{
		AuthorizationKey,
		RefreshTokenCookieName,
		AccessTokenCookieName,
	}
)