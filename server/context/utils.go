package context

import (
	"context"
	"net"
	"strings"

	"connectrpc.com/connect"
	goconnect "github.com/ralvarezdev/go-connect"
	goconnectserver "github.com/ralvarezdev/go-connect/server"
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
	// Try to get peer information from context first
	callInfo, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return "", goconnectserver.ErrFailedToGetCallInfo
	}
	
	// Get peer information
	peer := callInfo.Peer()
	if peer.Addr != "" {
		// Parse the address to extract IP (removes port if present)
		host, _, err := net.SplitHostPort(peer.Addr)
		if err != nil {
			// If SplitHostPort fails, it might be just an IP without port
			return peer.Addr, nil
		}
		return host, nil
	}

	// Get the request headers from the context
	headers, err := goconnectrequest.GetHeadersFromRequestContext(ctx)
	if err != nil {
		return "", err
	}
	
	// Check for X-Forwarded-For header
	if xForwardedFor := headers.Get(goconnect.XForwardedForKey); xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0]), nil
		}
	}

	// Check for X-Real-IP header
	if xRealIP := headers.Get(goconnect.XRealIPKey); xRealIP != "" {
		return xRealIP, nil
	}

	// Check for RemoteAddr
	if remoteAddr := headers.Get(goconnect.RemoteAddrKey); remoteAddr != "" {
		host, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			return remoteAddr, nil
		}
		return host, nil
	}

	return "", goconnect.ErrClientIPNotFound
}
