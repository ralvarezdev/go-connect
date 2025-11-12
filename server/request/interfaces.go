package request

import (
	"context"

	"connectrpc.com/connect"
)

type (
	// Injector is the interface that defines the methods to inject data into the request context
	Injector interface {
		CreateClientContextFromRequestContext(
			originalCtx context.Context, headers ...string,
		) (ctx context.Context, callInfo connect.CallInfo, err error)
	}
)
