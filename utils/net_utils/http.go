package net_utils

import (
	"net"
	"strings"
)

// RemoteInfo holds information about remote http client
type RemoteInfo struct {
	Addr, Host, Port string
}

// URLInfo - structure carrying information about current request and it's mapping to filesystem
type URLInfo struct {
	ScriptPath string
	PathInfo   string
	FilePath   string
}

func TellHostPort(host string, ssl bool) (server, port string, err error) {
	server, port, err = net.SplitHostPort(host)
	if err != nil {
		if addrerr, ok := err.(*net.AddrError); ok && strings.Contains(addrerr.Err, "missing port") {
			server = host
			if ssl {
				port = "443"
			} else {
				port = "80"
			}
			err = nil
		}
	}
	return server, port, err
}
