package context

import (
	"strings"
	"context"
	
	goconnect "github.com/ralvarezdev/go-connect"
	goconnectrequest "github.com/ralvarezdev/go-connect/server/request"
)

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