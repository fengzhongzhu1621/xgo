package client_transport

import (
	"context"
	"reflect"
	"sync"

	"github.com/fengzhongzhu1621/xgo/network/transport/options"
)

var (
	clientStreamTrans    = make(map[string]IClientStreamTransport)
	muxClientStreamTrans = sync.RWMutex{}
)

// IClientStreamTransport is the client stream transport interface.
// It's compatible with common RPC transport.
type IClientStreamTransport interface {
	// Send sends stream messages.
	Send(ctx context.Context, req []byte, opts ...options.RoundTripOption) error
	// Recv receives stream messages.
	Recv(ctx context.Context, opts ...options.RoundTripOption) ([]byte, error)
	// Init inits the stream.
	Init(ctx context.Context, opts ...options.RoundTripOption) error
	// Close closes stream transport, return connection to the resource pool.
	Close(ctx context.Context)
}

// RegisterClientStreamTransport registers a named IClientStreamTransport.
func RegisterClientStreamTransport(name string, t IClientStreamTransport) {
	tv := reflect.ValueOf(t)
	if t == nil || tv.Kind() == reflect.Ptr && tv.IsNil() {
		panic("transport: register nil client transport")
	}
	if name == "" {
		panic("transport: register empty name of client transport")
	}
	muxClientStreamTrans.Lock()
	clientStreamTrans[name] = t
	muxClientStreamTrans.Unlock()
}

// GetClientStreamTransport returns ClientStreamTransport by name.
func GetClientStreamTransport(name string) IClientStreamTransport {
	muxClientStreamTrans.RLock()
	t := clientStreamTrans[name]
	muxClientStreamTrans.RUnlock()
	return t
}
