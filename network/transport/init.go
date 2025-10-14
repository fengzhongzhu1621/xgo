package transport

import (
	"github.com/fengzhongzhu1621/xgo/network/transport/client_transport"
	"github.com/fengzhongzhu1621/xgo/network/transport/server_transport"
)

const transportName = "go-net"

func init() {
	client_transport.RegisterClientTransport(transportName, client_transport.DefaultClientTransport)
	client_transport.RegisterClientStreamTransport(transportName, client_transport.DefaultClientStreamTransport)
}

func init() {
	server_transport.RegisterServerTransport(transportName, server_transport.DefaultServerStreamTransport)
}
