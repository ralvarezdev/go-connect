# go-connect

Interceptors and helpers for building [ConnectRPC](https://connectrpc.com/) servers in Go: JWT authentication with automatic token refresh, error handling, header and cookie injection, client IP lookup and request validation.

## Installation

```bash
go get github.com/ralvarezdev/go-connect
```

Requires Go 1.25.3 or newer. Depends on `connectrpc.com/connect` and the author's `go-flags`, `go-grpc`, `go-jwt`, `go-reflect` and `go-validator`.

## Usage

Error-handling interceptor:

```go
import (
	"connectrpc.com/connect"
	goconnecterrorhandler "github.com/ralvarezdev/go-connect/server/interceptor/errorhandler"
)

errorHandler, err := goconnecterrorhandler.NewInterceptor(modeFlag, logger) // *goflagsmode.Flag, *slog.Logger
if err != nil {
	return err
}

path, handler := v1connect.NewAuthServiceHandler(
	service,
	connect.WithInterceptors(errorHandler.HandleError()),
)
mux.Handle(path, handler)
```

Authentication interceptor: `auth.NewInterceptor(modeFlag, validator, interceptions, options, logger)`, where `validator` is a `go-jwt` token validator, `interceptions` maps RPC procedure names to the expected token type (access or refresh), and `options.RefreshTokenFn` optionally refreshes an expired access token. Pass `interceptor.Authenticate()` to `connect.WithInterceptors`.

## Packages

- **`goconnect`** — header and cookie name constants (`AuthorizationKey`, `AccessTokenCookieName`, `RefreshTokenCookieName`, `AccessTokenKey`, `RefreshTokenKey`, `XForwardedForKey`, `RemoteAddrKey`, `XRealIPKey`), `AuthHeaders`, and errors (`ErrInvalidAuthorization`, `ErrMissingAuthorization`, `ErrClientIPNotFound`, `ErrInternalServerError`).
- **`server/interceptor/auth`** — `Interceptor`, `NewInterceptor`, `Authenticate()`, `Options`, `RefreshTokenFn`, `Authenticator`, `FindAuthorizationToken(header, customName, cookieName)`.
- **`server/interceptor/errorhandler`** — `Interceptor`, `NewInterceptor(modeFlag, logger)`, `HandleError()`, `ErrorHandler`.
- **`server/request`** — `Injector` and `DefaultInterceptor` (`CreateClientContextFromRequestContext`) to forward auth headers to downstream client calls, plus `GetHeadersFromRequestContext` and `SetHeadersToCallInfo`.
- **`server/response`** — `Injector` and `DefaultInterceptor` (`InjectTokens`, `InjectTokensFromContext`, `InjectHeadersFromCallInfo`) to write tokens and headers into responses.
- **`server/context`** — `GetClientIP`, `Get/SetCtxIssuedAccessToken`, `Get/SetCtxIssuedRefreshToken`.
- **`server/validator`** — `Service` and `DefaultService` (`NewService`) with `Email`, `Username`, `Birthdate`, `Password`, `CreateValidateFn` and `Validate`, built on `go-validator`.

There are no tests; linting is configured in `.golangci.yml`.

## License

GNU General Public License v3.0. See [LICENSE](LICENSE).
