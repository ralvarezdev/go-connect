package request

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	goconnectserver "github.com/ralvarezdev/go-connect/server"
)

// GetHeadersFromRequestContext tries to get the request headers from the context
//
// Parameters:
//
//   - ctx: The context to get the request headers from
//
// Returns:
//
// - http.Header: The request headers
// - error: If there was an error getting the request headers
func GetHeadersFromRequestContext(
	ctx context.Context,
) (http.Header, error) {
	// Get the call info from the context
	callInfo, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, goconnectserver.ErrFailedToGetCallInfo
	}

	// Get the request headers from the call info
	return callInfo.RequestHeader(), nil
}

// SetHeadersToCallInfo injects the specified headers from the headers into the call info
// 
// Parameters:
// 
//  - headers: The headers to inject
//  - callInfo: The call info to inject the headers into
//  - keys: The keys of the headers to inject
func SetHeadersToCallInfo(
	headers http.Header,
	callInfo connect.CallInfo,
	keys...string,
) {
	if callInfo == nil || headers == nil || len(keys) == 0 {
		return
	}
	
	// Get the request headers from the call info
	reqHeaders := callInfo.RequestHeader()
	
	// Iterate over the keys
	for _, key := range keys {
		// Check if the headers contain the key
		header := headers.Get(key)
		if header == "" {
			continue
		}

		// Inject the header
		reqHeaders.Set(key, header)
	}
}