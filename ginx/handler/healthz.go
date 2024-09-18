package handler

import (
	"fmt"
	"net/http"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/db/mysql"
	"github.com/fengzhongzhu1621/xgo/db/mysql/sqlxx"
	"github.com/fengzhongzhu1621/xgo/db/redisx"
	"github.com/gin-gonic/gin"
)

func NewHealthzHandleFunc(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. check database
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

		// 2. check redis
		var err error
		var addr string
		redisConfig, ok := cfg.RedisMap[redisx.ModeStandalone]
		if ok {
			addr = redisConfig.Addr
			err = redisx.TestConnection(&redisConfig)
		}

		redisConfig, ok = cfg.RedisMap[redisx.ModeSentinel]
		if ok {
			addr = redisConfig.SentinelAddr
			err = redisx.TestConnection(&redisConfig)
		}

		if err != nil {
			message := fmt.Sprintf("redis(mode=%s) connect fail: %s [addr=%s]", redisConfig.ID, err.Error(), addr)
			c.String(http.StatusInternalServerError, message)
			return
		}

		c.String(http.StatusOK, "ok")
	}
}
