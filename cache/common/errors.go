package common

import (
	"errors"
	"time"
)

// ErrNotExceptedTypeFromCache ...
var ErrNotExceptedTypeFromCache = errors.New("not expected type from cache")

const EmptyCacheExpiration = 5 * time.Second

// EmptyCache is a placeholder for the missing key
type EmptyCache struct {
	Err error
}
