package zookeeper

import "github.com/spf13/pflag"

type APIConfig struct {
	AddrPort    string
	RegDiscover string
	RegisterIP  string
	ExConfig    string
	Environment string
	Qps         int64
	Burst       int64
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{
		AddrPort:    "127.0.0.1:8081",
		RegDiscover: "",
		RegisterIP:  "",
		Qps:         1000,
		Burst:       2000,
	}
}

func (conf *APIConfig) AddFlags(fs *pflag.FlagSet, defaultAddrPort string) {
	fs.StringVar(&conf.AddrPort, "addrport", defaultAddrPort, "The ip address and port for the serve on")
	fs.StringVar(&conf.RegDiscover, "regdiscv", "", "hosts of register and discover server. e.g: 127.0.0.1:2181")
	fs.StringVar(&conf.ExConfig, "config", "", "The config path. e.g conf/api.conf")
	fs.StringVar(&conf.RegisterIP, "register-ip", "", "the ip address registered on zookeeper, it can be domain")
	fs.StringVar(&conf.Environment, "env", "", "the environment of the server, used for service discovery")
}
