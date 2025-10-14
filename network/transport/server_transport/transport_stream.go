package server_transport

import (
	"context"
	"reflect"
	"sync"
)

var (
	serverStreamTrans    = make(map[string]IServerStreamTransport)
	muxServerStreamTrans = sync.RWMutex{}
)

// IServerStreamTransport is the server stream transport interface.
// It's compatible with common RPC transport.
type IServerStreamTransport interface {
	// ServerTransport is used to keep compatibility with common RPC transport.
	IServerTransport
	// Send sends messages.
	Send(ctx context.Context, req []byte) error
	// Close is called when server encounters an error and cleans up.
	Close(ctx context.Context)
}

// RegisterServerStreamTransport Registers a named ServerStreamTransport.
func RegisterServerStreamTransport(name string, t IServerStreamTransport) {
	tv := reflect.ValueOf(t)
	if t == nil || tv.Kind() == reflect.Ptr && tv.IsNil() {
		panic("transport: register nil server transport")
	}
	if name == "" {
		panic("transport: register empty name of server transport")
	}
	muxServerStreamTrans.Lock()
	serverStreamTrans[name] = t
	muxServerStreamTrans.Unlock()

}

// GetServerStreamTransport returns ServerStreamTransport by name.
func GetServerStreamTransport(name string) IServerStreamTransport {
	muxServerStreamTrans.RLock()
	t := serverStreamTrans[name]
	muxServerStreamTrans.RUnlock()
	return t
}
