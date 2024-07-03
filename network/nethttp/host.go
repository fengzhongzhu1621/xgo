package nethttp

import (
	"net"
	"strings"
)

// 获得主机和端口，如果不存在则获取默认值.
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
