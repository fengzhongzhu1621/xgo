package nethttp

import (
	"net"
)

// RemoteInfo holds information about remote http client.
type RemoteInfo struct {
	Addr, Host, Port string
}

// URLInfo - structure carrying information about current request and it's mapping to filesystem.
type URLInfo struct {
	ScriptPath string
	PathInfo   string
	FilePath   string
}

// GetRemoteInfo creates RemoteInfo structure and fills its fields appropriately.
func GetRemoteInfo(remote string, doLookup bool) (*RemoteInfo, error) {
	// 获得IP和端口
	addr, port, err := net.SplitHostPort(remote)
	if err != nil {
		return nil, err
	}

	var host string
	if doLookup {
		// 用于根据给定的 IP 地址查找其对应的域名（反向 DNS 查询）。
		// 这个函数返回与给定 IP 地址关联的所有主机名。
		hosts, err := net.LookupAddr(addr)
		if err != nil || len(hosts) == 0 {
			host = addr
		} else {
			host = hosts[0]
		}
	} else {
		host = addr
	}

	return &RemoteInfo{Addr: addr, Host: host, Port: port}, nil
}
