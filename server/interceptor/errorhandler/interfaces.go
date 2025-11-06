package errorhandler

import (
	"connectrpc.com/connect"
)

type (
	// ErrorHandler interface
	ErrorHandler interface {
		HandleError() connect.UnaryInterceptorFunc
	}
)
