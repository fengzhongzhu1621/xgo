package redisx

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/db/redisx/client"
	redis "github.com/fengzhongzhu1621/xgo/db/redisx/client"
	log "github.com/sirupsen/logrus"
)

// InitRedis 根据全局变量初始化 redis 连接
func InitRedis() {
	// 从全局配置获取 redis 配置
	globalConfig := config.GetGlobalConfig()
	redisMap := globalConfig.RedisMap
	// 获取 redis 集群的类型
	standaloneConfig, isStandalone := redisMap[redis.ModeStandalone]
	sentinelConfig, isSentinel := redisMap[redis.ModeSentinel]

	if !(isStandalone || isSentinel) {
		panic("redis id=standalone or id=sentinel should be configured")
	}

	if isSentinel && isStandalone {
		log.Info("redis both id=standalone and id=sentinel configured, will use sentinel")

		delete(globalConfig.RedisMap, redis.ModeStandalone)
		isStandalone = false
	}

	if isSentinel {
		if sentinelConfig.MasterName == "" {
			panic("redis id=sentinel, the `masterName` required")
		}
		log.Info("init Redis mode=`sentinel`")
		client.InitRedisClient(&sentinelConfig)
	}

	if isStandalone {
		log.Info("init Redis mode=`standalone`")
		client.InitRedisClient(&standaloneConfig)
	}

	log.Info("init Redis success")
}

func init() {
	InitRedis()
}
