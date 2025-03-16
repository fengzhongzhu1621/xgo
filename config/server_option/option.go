package zookeeper

import (
	"github.com/fengzhongzhu1621/xgo/config/flagutils"
	"github.com/spf13/pflag"
)

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

type ServerOption struct {
	ServConf         *APIConfig
	DeploymentMethod flagutils.DeploymentMethod
}

func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf:         NewAPIConfig(),
		DeploymentMethod: flagutils.OpenSourceDeployment,
	}
	return &s
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	s.ServConf.AddFlags(fs, "127.0.0.1:50001")
	fs.Var(flagutils.EnableAuthFlag, "enable-auth", "The auth center enable status, true for enabled, false for disabled")
	fs.Var(&s.DeploymentMethod, "deployment-method", "The deployment method, supported value: open_source, blueking")

}

// SetServerInfo Information about the current process in service governance
func SetServerInfo(srvInfo *ServerInfo) {
	server = srvInfo
}

// GetServerInfo Information about the current process in service governance
func GetServerInfo() *ServerInfo {
	return server
}
