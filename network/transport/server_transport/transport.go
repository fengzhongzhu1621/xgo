package server_transport

import (
	"context"
	"reflect"
	"sync"

	"github.com/fengzhongzhu1621/xgo/network/transport/options"
)

var (
	svrTrans    = make(map[string]IServerTransport)
	muxSvrTrans = sync.RWMutex{}
)

// ServerTransport defines the server transport layer interface.
type IServerTransport interface {
	ListenAndServe(ctx context.Context, opts ...options.ListenServeOption) error
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// RegisterServerTransport register a ServerTransport.
func RegisterServerTransport(name string, t IServerTransport) {
	tv := reflect.ValueOf(t)
	if t == nil || tv.Kind() == reflect.Ptr && tv.IsNil() {
		panic("transport: register nil server transport")
	}
	if name == "" {
		panic("transport: register empty name of server transport")
	}
	muxSvrTrans.Lock()
	svrTrans[name] = t
	muxSvrTrans.Unlock()
}

// GetServerTransport gets the ServerTransport.
func GetServerTransport(name string) IServerTransport {
	muxSvrTrans.RLock()
	t := svrTrans[name]
	muxSvrTrans.RUnlock()
	return t
}
