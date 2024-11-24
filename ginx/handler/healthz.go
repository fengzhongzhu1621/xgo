package handler

import (
	"fmt"
	"github.com/fengzhongzhu1621/xgo/ginx/serializer"
	"net/http"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/db/mysql"
	"github.com/fengzhongzhu1621/xgo/db/mysql/sqlxx"
	redis "github.com/fengzhongzhu1621/xgo/db/redisx/client"
	"github.com/gin-gonic/gin"
)

func NewHealthzHandleFunc(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check database
		defaultDBConfig := cfg.DatabaseMap["default"]
		dbConfigs := []mysql.Database{defaultDBConfig}
		for _, dbConfig := range dbConfigs {
			dbConfig := dbConfig
			// reset the options for check
			dbConfig.MaxIdleConns = 1
			dbConfig.MaxOpenConns = 1
			dbConfig.ConnMaxLifetimeSecond = 60

			// 测试数据库连接
			err := sqlxx.TestConnection(&dbConfig)
			if err != nil {
				message := fmt.Sprintf("db connect fail: %s [id=%s host=%s port=%d]",
					err.Error(), dbConfig.ID, dbConfig.Host, dbConfig.Port)
				c.String(http.StatusInternalServerError, message)
				return
			}
		}

		// check redis
		var err error
		var addr string
		redisConfig, ok := cfg.RedisMap[redis.ModeStandalone]
		if ok {
			addr = redisConfig.Addr
			err = redis.TestConnection(&redisConfig)
		}

		redisConfig, ok = cfg.RedisMap[redis.ModeSentinel]
		if ok {
			addr = redisConfig.SentinelAddr
			err = redis.TestConnection(&redisConfig)
		}

		if err != nil {
			message := fmt.Sprintf("redis(mode=%s) connect fail: %s [addr=%s]", redisConfig.Type, err.Error(), addr)
			c.String(http.StatusInternalServerError, message)
			return
		}

		c.JSON(http.StatusOK, serializer.HealthResponse{Healthy: true})
	}
}
