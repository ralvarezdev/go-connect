package errorhandler

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"

	goconnect "github.com/ralvarezdev/go-connect"
	goflags "github.com/ralvarezdev/go-flags"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Interceptor is the interceptor for the error handler
	Interceptor struct {
		modeFlag *goflagsmode.Flag
		logger *slog.Logger
	}
)

// NewInterceptor creates a new error handler interceptor
//
// Parameters:
//
//   - modeFlag: the application mode flag
//   - logger: the logger to use (can be nil)
//
// Returns:
//
//   - *Interceptor: the interceptor
//   - error: if there was an error creating the interceptor
func NewInterceptor(modeFlag *goflagsmode.Flag, logger *slog.Logger) (*Interceptor, error) {
	// Check if the mode flag is nil
	if modeFlag == nil {
	 	return nil, goflags.ErrNilFlag
	}
	
	// Create the logger for the interceptor
	if logger != nil {
		logger = logger.With(
			slog.String("grpc_client_interceptor", "error_handler"),
		)
	}

	return &Interceptor{
		modeFlag: modeFlag,
		logger: logger,
	}, nil
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
					stack := debug.Stack()
					if i.logger != nil {
						i.logger.Error(
							"Panic recovered",
							slog.Any("method", req.Spec().Procedure),
							slog.Any("error", r),
							slog.String("stack_trace", string(stack)),
						)
					}

					// Check if the application is in production mode
					if i.modeFlag.IsProd() {
						// Set the error to internal server error
						err = connect.NewError(connect.CodeInternal, goconnect.ErrInternalServerError)
					} else {
						// Set the error to the recovered panic
						err = connect.NewError(
							connect.CodeInternal,
							fmt.Errorf(
								"Panic: %v\nStack Trace:\n%s",
								r,
								string(stack),
							),
						)
					}
				}
			}()
			return next(ctx, req)
		}
	}
}
