package errorhandler

import (
	"context"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gogrpc "github.com/ralvarezdev/go-grpc"
)

type (
	// Interceptor is the interceptor for the error handler
	Interceptor struct {
		logger *slog.Logger
	}
)

// NewInterceptor creates a new error handler interceptor
//
// Parameters:
//
//   - logger: the logger to use (can be nil)
//
// Returns:
//
//   - *Interceptor: the interceptor
func NewInterceptor(logger *slog.Logger) *Interceptor {
	if logger != nil {
		logger = logger.With(
			slog.String("grpc_client_interceptor", "error_handler"),
		)
	}

	return &Interceptor{
		logger: logger,
	}
}

// HandleError returns the error handler interceptor
//
// Returns:
//
//   - connect.UnaryInterceptorFunc: the error handler interceptor
func (i Interceptor) HandleError() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (res connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic
					if i.logger != nil {
						i.logger.Error(
							"Panic recovered",
							slog.Any("method", req.Spec().Procedure),
							slog.Any("error", r),
							slog.String("stack_trace", string(debug.Stack())),
						)
					}

					// Set the error to internal server error
					err = status.Error(codes.Internal, gogrpc.InternalServerError)
				}
			}()
			return next(ctx, req)
		}
	}
}
