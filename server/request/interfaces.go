package request

import "context"

type (
	// Injector is the interface that defines the methods to inject data into the request context
	Injector interface {
		CreateClientContextFromRequestContext(
			ctx context.Context, headers ...string,
		) (context.Context, error)
	}
)
