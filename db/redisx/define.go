package redisx

type Redis struct {
	ID           string // redis 的集群类型
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
}
