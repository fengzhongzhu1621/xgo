package config

import (
	"github.com/fengzhongzhu1621/xgo/db/kafkax"
	"github.com/fengzhongzhu1621/xgo/db/mysql"
	"github.com/fengzhongzhu1621/xgo/db/redisx"
	"github.com/gin-gonic/gin"
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

	// 数据库配置
	Databases   []mysql.Database
	DatabaseMap map[string]mysql.Database

	// redis 配置
	Redis    []redisx.Redis
	RedisMap map[string]redisx.Redis

	// kafka 配置
	Kafka    []kafkax.Kafka
	KafkaMap map[string]kafkax.Kafka

	// pprof
	PprofAccounts gin.Accounts // 认证用户

	// 日志
	Logger Logger

	// 版本
	RootDir string
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
