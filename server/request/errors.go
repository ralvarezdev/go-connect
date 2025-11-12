package request

import (
	"errors"
)

var (
	ErrNilInjector = errors.New("nil request injector")
)
