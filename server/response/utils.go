package response

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	goconnectserver "github.com/ralvarezdev/go-connect/server"
)

// GetHeadersFromRequestContext tries to get the response headers from the context
//
// Parameters:
//
// - ctx: The context to get the response headers from
//
// Returns:
//
// - http.Header: The response headers
// - error: If there was an error getting the response headers
func GetHeadersFromRequestContext(
	ctx context.Context,
) (http.Header, error) {
	// Get the call info from the context
	callInfo, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, goconnectserver.ErrFailedToGetCallInfo
	}

	// Get the response headers from the call info
	return callInfo.ResponseHeader(), nil
}

// SetHeadersFromCallInfo injects the headers from the call info to the headers
//
// Parameters:
// 
// - respHeader: The headers to inject the headers to
// - callInfo: The call info to get the headers from
func SetHeadersFromCallInfo(
	headers http.Header,
	callInfo connect.CallInfo,
) {
	// Check if the headers or call info are nil
	if headers == nil || callInfo == nil {
		return
	}
	
	// Inject the headers from the call info into the headers
	for key, values := range callInfo.ResponseHeader() {
		// Add all values for the key
		for _, value := range values {
			headers.Add(key, value)
		}
	}
}