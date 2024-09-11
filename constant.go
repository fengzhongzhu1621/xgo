package xgo

import (
	"math"
	"time"
)

const (
	RequestIDKey       = "request_id"
	RequestIDHeaderKey = "X-Request-Id"
)

var (
	NeverExpires = time.Unix(math.MaxInt64, 0)
)
