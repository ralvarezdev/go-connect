package auth

import (
	"net/http"
	"strings"

	goconnect "github.com/ralvarezdev/go-connect"
	gojwt "github.com/ralvarezdev/go-jwt"
)

// FindAuthorizationToken extracts the authorization token from the HTTP headers.
//
// Parameters:
//
//  - header: the HTTP headers to extract the token from
//  - customName: optional custom header name to check for the token
//  - cookieName: the name of the cookie to check for the token
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
