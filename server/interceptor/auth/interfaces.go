package auth

import (
	"connectrpc.com/connect"
)

type (
	// Authenticator interface
	Authenticator interface {
		Authenticate() connect.UnaryInterceptorFunc
	}
)
