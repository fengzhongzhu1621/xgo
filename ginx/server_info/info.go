package server_info

import (
	"fmt"
	"os"

	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/network/ip"
	"github.com/fengzhongzhu1621/xgo/version"
	"github.com/rs/xid"
)

var server *ServerInfo

type ServerInfo struct {
	IP         string `json:"ip"`
	Port       uint   `json:"port"`
	RegisterIP string `json:"registerip"`
	HostName   string `json:"hostname"`
	Scheme     string `json:"scheme"`
	Version    string `json:"version"`
	Pid        int    `json:"pid"`
	// UUID is used to distinguish which service is master in zookeeper
	UUID string `json:"uuid"`
	// Environment is the server's environment, servers can only discover other servers in the same environment
	Environment string `json:"env"`
}

// SetServerInfo Information about the current process in service governance
func SetServerInfo(srvInfo *ServerInfo) {
	server = srvInfo
}

// GetServerInfo Information about the current process in service governance
func GetServerInfo() *ServerInfo {
	return server
}

// NewServerInfo new a ServerInfo object，配置是从命令行获取的
func NewServerInfo(conf *server_option.APIConfig) (*ServerInfo, error) {
	ipValue, err := ip.GetAddress(conf.AddrPort)
	if err != nil {
		return nil, err
	}

	port, err := ip.GetPort(conf.AddrPort)
	if err != nil {
		return nil, err
	}

	registerIP := conf.RegisterIP
	// if no registerIP is set, default to be the ip
	if registerIP == "" {
		registerIP = ipValue
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &ServerInfo{
		IP:          ipValue,
		Port:        port,
		RegisterIP:  registerIP,
		HostName:    hostname,
		Scheme:      "http",
		Version:     version.GetVersion(),
		Pid:         os.Getpid(),
		UUID:        xid.New().String(),
		Environment: conf.Environment,
	}
	return info, nil
}

func (s *ServerInfo) RegisterAddress() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.RegisterIP, s.Port)
}

func (s *ServerInfo) Instance() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
