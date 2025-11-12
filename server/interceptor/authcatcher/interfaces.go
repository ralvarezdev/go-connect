package auth

import (
	"connectrpc.com/connect"
)

type (
	// AuthenticatorCatcher interface
	AuthenticatorCatcher interface {
		Catch() connect.UnaryInterceptorFunc
	}
)