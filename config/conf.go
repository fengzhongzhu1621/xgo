package config

type Database struct {
	ID       string
	Host     string
	Port     int
	User     string
	Password string
	Name     string

	MaxOpenConns          int
	MaxIdleConns          int
	ConnMaxLifetimeSecond int
}

type Redis struct {
	ID           string
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

type Kafka struct {
	Id       string `yaml:"id"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Topic    string `yaml:"topic"`
	Version  string `yaml:"version"`
	LogPath  string `yaml:"log_path"`
}

type Config struct {
	// 数据库配置
	Databases   []Database
	DatabaseMap map[string]Database

	Redis    []Redis
	RedisMap map[string]Redis

	Kafka    []Kafka
	KafkaMap map[string]Kafka
}
