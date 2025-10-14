package client_transport

import (
	"context"
	"reflect"
	"sync"

	"github.com/fengzhongzhu1621/xgo/network/transport/options"
)

var (
	clientTrans    = make(map[string]IClientTransport)
	muxClientTrans = sync.RWMutex{}
)

// ClientTransport defines the client transport layer interface.
type IClientTransport interface {
	RoundTrip(ctx context.Context, req []byte, opts ...options.RoundTripOption) (rsp []byte, err error)
}

// RegisterClientTransport register a ClientTransport.
func RegisterClientTransport(name string, t IClientTransport) {
	tv := reflect.ValueOf(t)
	if t == nil || tv.Kind() == reflect.Ptr && tv.IsNil() {
		panic("transport: register nil client transport")
	}
	if name == "" {
		panic("transport: register empty name of client transport")
	}
	muxClientTrans.Lock()
	clientTrans[name] = t
	muxClientTrans.Unlock()
}

// GetClientTransport gets the ClientTransport.
func GetClientTransport(name string) IClientTransport {
	muxClientTrans.RLock()
	t := clientTrans[name]
	muxClientTrans.RUnlock()
	return t
}
