package config

import (
	"github.com/fengzhongzhu1621/xgo/config/flagutils"
	kafka "github.com/fengzhongzhu1621/xgo/db/kafka"
	"github.com/fengzhongzhu1621/xgo/db/mysql"
	redis "github.com/fengzhongzhu1621/xgo/db/redis/client"
	"github.com/spf13/viper"
)

type LogConfig struct {
	Level    string
	Writer   string
	Settings map[string]string // 日志详细配置
}

type Logger struct {
	System LogConfig
	API    LogConfig
	Web    LogConfig
}

type ServiceLogConfig struct {
	Level string
	Dir   string
}

type Config struct {
	// 调试开关
	Debug bool

	EnableSwagger bool

	Server Server

	Sentry Sentry

	// 数据库配置
	Databases   []mysql.Database
	DatabaseMap map[string]mysql.Database

	// redis 配置
	Redis    []redis.Redis
	RedisMap map[string]redis.Redis

	// kafka 配置
	Kafka    []kafka.Kafka
	KafkaMap map[string]kafka.Kafka

	// pprof
	PProf PProf `yaml:"pprof"`

	// 日志
	Logger Logger

	// 版本
	RootDir string

	// 认证
	Auth Auth

	Cryptos map[string]*Crypto

	SuperAppCode []string
}

type Server struct {
	Host string
	Port int

	GraceTimeout int64

	// 服务器在读取请求体时的最大持续时间。如果超过这个时间，服务器将中断读取并返回错误。
	ReadTimeout int
	// 服务器在写入响应体时的最大持续时间。如果超过这个时间，服务器将中断写入并返回错误。
	WriteTimeout int
	// 服务器在关闭连接之前等待下一个请求的最大时间。这对于管理服务器资源很有用，特别是在高并发场景下。
	IdleTimeout int

	// 如果是 https，必须设置 TlsCertFile 和 TlsKeyFile
	Mode        string
	TlsCertFile string
	TlsKeyFile  string
}

type Sentry struct {
	Enable bool
	DSN    string
}

type PProf struct {
	// 认证用户
	Account map[string]string
}

type Auth struct {
	BearerToken string
	JwtToken    string
}

type Crypto struct {
	ID  string
	Key string
}

type WebServerConfig struct {
	Site                      Site
	Session                   Session
	Redis                     redis.Config
	Version                   string
	AgentAppUrl               string
	LoginUrl                  string
	LoginVersion              string
	ConfigMap                 map[string]string
	AuthCenter                AppInfo
	DisableOperationStatistic bool
	DeploymentMethod          flagutils.DeploymentMethod
	EnableNotification        bool
}

type Site struct {
	AccountUrl      string
	DomainUrl       string
	HttpsDomainUrl  string
	HtmlRoot        string
	ResourcesPath   string
	BkLoginUrl      string
	BkHttpsLoginUrl string
	AppCode         string
	CheckUrl        string
	// available value: internal, iam
	AuthScheme string
	// available value: off, on
	FullTextSearch string
	PaasDomainUrl  string
	// BkDomain it is used to set the unified international language of Blue Whale.
	// this variable is returned to the front-end through configuration in the background.
	// the corresponding front-end variable is: cookieDomain.
	BkDomain string
	// BkComponentApiUrl is the blueking component api url, and is also the esb domain url
	// the corresponding front-end variable is: componentApiUrl.
	BkComponentApiUrl string
	HelpDocUrl        string
	// BkSharedResUrl is the blueking shared resource url
	BkSharedResUrl string
}

type Session struct {
	Name            string
	DefaultLanguage string
	MultipleOwner   string
}

type AppInfo struct {
	AppCode string `json:"appCode"`
	URL     string `json:"url"`
}

// Load 将配置文件转换为全局结构体对象
func Load(v *viper.Viper) (*Config, error) {
	var cfg Config

	// 配置文件转换为全局结构体对象
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
