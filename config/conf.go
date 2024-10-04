package config

import (
	"github.com/fengzhongzhu1621/xgo/db/kafkax"
	"github.com/fengzhongzhu1621/xgo/db/mysql"
	redis "github.com/fengzhongzhu1621/xgo/db/redisx/client"
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

type Config struct {
	// 调试开关
	Debug bool

	Sentry Sentry

	// 数据库配置
	Databases   []mysql.Database
	DatabaseMap map[string]mysql.Database

	// redis 配置
	Redis    []redis.Redis
	RedisMap map[string]redis.Redis

	// kafka 配置
	Kafka    []kafkax.Kafka
	KafkaMap map[string]kafkax.Kafka

	// pprof
	PProf PProf `yaml:"pprof"`

	// 日志
	Logger Logger

	// 版本
	RootDir string
}

type Sentry struct {
	Enable bool
	DSN    string
}

type PProf struct {
	// 认证用户
	Account map[string]string
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
