package xgo

import (
	"math"
	"time"
)

var NeverExpires = time.Unix(math.MaxInt64, 0)

// InfiniteEndID represent infinity for end id of id rule info
const InfiniteEndID int64 = -1
