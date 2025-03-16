package server_option

import (
	"github.com/fengzhongzhu1621/xgo/config/flagutils"
	"github.com/spf13/pflag"
)

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
