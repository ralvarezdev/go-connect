package headers

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