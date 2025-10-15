package connpool

import (
	"errors"
	"net"
)

// connection pool error message.
var (
	ErrPoolLimit  = errors.New("connection pool limit")  // ErrPoolLimit number of connections exceeds the limit error.
	ErrPoolClosed = errors.New("connection pool closed") // ErrPoolClosed connection pool closed error.
	ErrConnClosed = errors.New("conn closed")            // ErrConnClosed connection closed.
	ErrNoDeadline = errors.New("dial no deadline")       // ErrNoDeadline has no deadline set.
	ErrConnInPool = errors.New("conn already in pool")   // ErrNoDeadline has no deadline set.
)

// Pool is the interface that specifies client connection pool options.
// Compared with Pool, Pool directly uses the GetOptions data structure for function input parameters.
// Compared with function option input parameter mode, it can reduce memory escape and improve calling performance.
type IPool interface {
	Get(network string, address string, opt GetOptions) (net.Conn, error)
}

// HealthChecker idle connection health check function.
// The function supports quick check and comprehensive check.
// Quick check is called when an idle connection is obtained,
// and only checks whether the connection status is abnormal.
// The function returns true to indicate that the connection is available normally.
type IHealthChecker func(pc *PoolConn, isFast bool) bool
