package transport

import (
	"context"
	"net"
	"reflect"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/transport/client_transport"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/network/transport/server_transport"
)

// contextKey is the context key.
type contextKey struct {
	name string
}

var (
	// LocalAddrContextKey is the local address context key.
	LocalAddrContextKey = &contextKey{"local-addr"}
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var (
	// RemoteAddrContextKey is the remote address context key.
	RemoteAddrContextKey = &contextKey{"remote-addr"}
)

// RemoteAddrFromContext gets remote address from context.
func RemoteAddrFromContext(ctx context.Context) net.Addr {
	addr, ok := ctx.Value(RemoteAddrContextKey).(net.Addr)
	if !ok {
		return nil
	}
	return addr
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var framerBuilders = make(map[string]codec.IFramerBuilder)

// FramerBuilder is the alias of codec.IFramerBuilder.
type FramerBuilder = codec.IFramerBuilder

// Framer is the alias of codec.IFramer.
type Framer = codec.IFramer

// RegisterFramerBuilder register a codec.IFramerBuilder.
func RegisterFramerBuilder(name string, fb codec.IFramerBuilder) {
	fbv := reflect.ValueOf(fb)
	if fb == nil || fbv.Kind() == reflect.Ptr && fbv.IsNil() {
		panic("transport: register framerBuilders nil")
	}
	if name == "" {
		panic("transport: register framerBuilders name empty")
	}
	framerBuilders[name] = fb
}

// GetFramerBuilder gets the FramerBuilder by name.
func GetFramerBuilder(name string) codec.IFramerBuilder {
	return framerBuilders[name]
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ListenAndServe wraps and starts the default server transport.
func ListenAndServe(opts ...options.ListenServeOption) error {
	return server_transport.DefaultServerTransport.ListenAndServe(context.Background(), opts...)
}

// RoundTrip wraps and starts the default client transport.
func RoundTrip(ctx context.Context, req []byte, opts ...options.RoundTripOption) ([]byte, error) {
	return client_transport.DefaultClientTransport.RoundTrip(ctx, req, opts...)
}
