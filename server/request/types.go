package request

import (
	"context"

	"connectrpc.com/connect"
)

type (
	// DefaultInterceptor is the default interceptor
	DefaultInterceptor struct{}
)

// CreateClientContextFromRequestContext creates a client context from the request context, injecting specified headers
//
// Parameters:
//
// - originalCtx: the request context
// - headers: the headers to inject
//
// Returns:
//
// - context.Context: the client context
// - error: if there was an error creating the client context
func (i DefaultInterceptor) CreateClientContextFromRequestContext(
	originalCtx context.Context, headers ...string,
) (ctx context.Context,  callInfo connect.CallInfo, err error) {
	// Get the call info from the context
	reqHeaders, err := GetHeadersFromRequestContext(originalCtx)
	if err != nil {
		return nil, nil, err
	}
	
	// Create the client context
	clientCtx, callInfo := connect.NewClientContext(context.Background())
	
	// Inject the headers into the client context
	for _, header := range headers {
		// Get the request header value
		reqHeader := reqHeaders.Get(header)
		
		// Inject the header
		callInfo.RequestHeader().Set(header, reqHeader)
	}
	
	return clientCtx, nil, err	
}
