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

/**
 * 获得主机和端口，如果不存在则获取默认值
 */
func TellHostPort(host string, ssl bool) (server, port string, err error) {
	// 根据字符串获得IP和端口
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
