package client

import "fmt"

const (
	NameCache = "cache"
	NameMQ    = "mq"
)

type Redis struct {
	Type         string // redis 的集群类型
	Addr         string
	Password     string
	DB           int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
	PoolSize     int
	MinIdleConns int
	ChannelKey   string
	SupportBRPOP bool

	// mode=sentinel required
	SentinelAddr     string
	MasterName       string
	SentinelPassword string

	debugMode bool
}

type RedisConfig struct {
	Username string
	Host     string
	Port     int
	Password string
	DB       int
}

func (cfg *RedisConfig) DSN() string {
	return fmt.Sprintf(
		"redis://%s:%s@%s:%d/%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
}
