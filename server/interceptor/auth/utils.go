package auth

import (
	"context"
	"net/http"
	"strings"

	gojwt "github.com/ralvarezdev/go-jwt"

	goconnect "github.com/ralvarezdev/go-connect"
	goconnectrequest "github.com/ralvarezdev/go-connect/server/request"
)

// FindAuthorizationToken extracts the authorization token from the HTTP headers.
//
// Parameters:
//
//   - header: the HTTP headers to extract the token from
//   - customName: optional custom header name to check for the token
//   - cookieName: the name of the cookie to check for the token
//
// Returns:
//
// - string: the extracted authorization token, or an empty string if not found
// - error: if there was an error during extraction
func FindAuthorizationToken(header http.Header, customName, cookieName *string) (token string, err error) {
	// Check for the Authorization header
	authHeader := header.Get(goconnect.AuthorizationKey)
	if authHeader == "" && customName != nil && *customName != "" {
		// Check for the custom header if provided
		authHeader = header.Get(*customName)
	}

	// If Authorization header is found, process it
	if authHeader != "" {
		// Split the authorization value by space
		authFields := strings.Split(authHeader, " ")

		// Check if the authorization value is valid
		if len(authFields) != 2 || authFields[0] != gojwt.BearerPrefix {
			return "", goconnect.ErrInvalidAuthorization
		}

		return authFields[1], nil
	}

	// Check if cookieName is provided
	if cookieName == nil || *cookieName == "" {
		return "", goconnect.ErrMissingAuthorization
	}

	// Check for the token in the specified cookie
	authCookie := header.Get(*cookieName)
	if authCookie == "" {
		return "", goconnect.ErrMissingAuthorization
	}
	return authCookie, nil
}

// GetClientIP retrieves the client IP address from the HTTP headers.
// 
// Parameters:
// 
// - ctx: The HTTP headers to extract the client IP from
// 
// Returns:
// 
// - string: The client IP address, or an empty string if not found
// - error: If there was an error during extraction
func GetClientIP(ctx context.Context) (string, error) {
	// Get the HTTP headers from the context
	headers, err := goconnectrequest.GetHeadersFromRequestContext(ctx)
	if err != nil {
		return "", err
	}

	// Check for the X-Forwarded-For header
	xForwardedFor := headers.Get(goconnect.XForwardedForKey)
	if xForwardedFor != "" {
		// The X-Forwarded-For header can contain multiple IPs, take the first one
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0]), nil
		}
	}

	// Fallback to RemoteAddr header if X-Forwarded-For is not present
	remoteAddr := headers.Get(goconnect.RemoteAddrKey)
	if remoteAddr != "" {
		return remoteAddr, nil
	}

	return "", goconnect.ErrClientIPNotFound
}