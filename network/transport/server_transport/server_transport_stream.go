package server_transport

import (
	"context"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/ip"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/xerror"
)

// serverStreamTransport implements ServerStreamTransport and keeps backward compatibility with the
// original serverTransport.
type serverStreamTransport struct {
	// Keep backward compatibility with original serverTransport.
	serverTransport
}

// NewServerStreamTransport creates a new ServerTransport, which is wrapped in serverStreamTransport
// as the return ServerStreamTransport interface.
func NewServerStreamTransport(opt ...options.ServerTransportOption) IServerStreamTransport {
	s := newServerTransport(opt...)
	return &serverStreamTransport{s}
}

// DefaultServerStreamTransport is the default ServerStreamTransport.
var DefaultServerStreamTransport = NewServerStreamTransport()

// ListenAndServe implements ServerTransport.
// To be compatible with common RPC and stream RPC, we use serverTransport.ListenAndServe function.
func (st *serverStreamTransport) ListenAndServe(ctx context.Context, opts ...options.ListenServeOption) error {
	return st.serverTransport.ListenAndServe(ctx, opts...)
}

// Send is the method to send stream messages.
func (st *serverStreamTransport) Send(ctx context.Context, req []byte) error {
	msg := codec.Message(ctx)
	raddr := msg.RemoteAddr()
	laddr := msg.LocalAddr()
	if raddr == nil || laddr == nil {
		return xerror.NewFrameError(xerror.RetServerSystemErr,
			fmt.Sprintf("Address is invalid, local: %s, remote: %s", laddr, raddr))
	}
	key := ip.AddrToKey(laddr, raddr)
	st.serverTransport.m.RLock()
	tc, ok := st.serverTransport.addrToConn[key]
	st.serverTransport.m.RUnlock()
	if ok && tc != nil {
		if _, err := tc.rwc.Write(req); err != nil {
			tc.close()
			st.Close(ctx)
			return err
		}
		return nil
	}
	return xerror.NewFrameError(xerror.RetServerSystemErr, "Can't find conn by addr")
}

// Close closes ServerStreamTransport, it also cleans up cached connections.
func (st *serverStreamTransport) Close(ctx context.Context) {
	msg := codec.Message(ctx)
	key := ip.AddrToKey(msg.LocalAddr(), msg.RemoteAddr())
	st.m.Lock()
	delete(st.serverTransport.addrToConn, key)
	st.m.Unlock()
}
