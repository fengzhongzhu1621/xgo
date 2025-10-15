package multiplexed

import (
	"context"
	"errors"
	"io"
	"net"
)

var (
	// ErrFrameParserNil indicates that frame parse is nil.
	ErrFrameParserNil = errors.New("frame parser is nil")
	// ErrRecvQueueFull receive queue full.
	ErrRecvQueueFull = errors.New("virtual connection's recv queue is full")
	// ErrSendQueueFull send queue is full.
	ErrSendQueueFull = errors.New("connection's send queue is full")
	// ErrChanClose connection is closed.
	ErrChanClose = errors.New("unexpected recv chan close")
	// ErrAssertFail type assert fail.
	ErrAssertFail = errors.New("type assert fail")
	// ErrDupRequestID duplicated request id.
	ErrDupRequestID = errors.New("duplicated Request ID")
	// ErrInitPoolFail failed to initialize connection.
	ErrInitPoolFail = errors.New("init pool for specific node fail")
	// ErrWriteNotFinished write operation is not completed.
	ErrWriteNotFinished = errors.New("write not finished")
	// ErrNetworkNotSupport does not support network type.
	ErrNetworkNotSupport = errors.New("network not support")
	// ErrConnectionsHaveBeenExpelled denotes that the connections to a certain ip:port have been expelled.
	ErrConnectionsHaveBeenExpelled = errors.New("connections have been expelled")
)

// FrameParser is the interface to parse a single frame.
type IFrameParser interface {
	// Parse parses vid and frame from io.ReadCloser. rc.Close must be called before Parse return.
	Parse(rc io.Reader) (vid uint32, buf []byte, err error)
}

// Pool is a connection pool for multiplexing.
type IPool interface {
	// GetMuxConn gets a multiplexing connection to the address on named network.
	GetMuxConn(ctx context.Context, network string, address string, opts GetOptions) (IMuxConn, error)
}

// IMuxConn is virtual connection multiplexing on a real connection.
type IMuxConn interface {
	// Write writes data to the connection.
	Write([]byte) error

	// Read reads a packet from connection.
	Read() ([]byte, error)

	// LocalAddr returns the local network address, if known.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address, if known.
	RemoteAddr() net.Addr

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close()
}
