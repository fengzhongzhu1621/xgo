package transport

import (
	stdhttp "net/http"
)

// ServerTransport is the http transport layer.
type ServerTransport struct {
	newServer func() *stdhttp.Server
	reusePort bool
	enableH2C bool
}
