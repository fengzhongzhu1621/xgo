package handler

import (
	"fmt"
	"net/http"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/db/mysql/sqlxx"
	"github.com/gin-gonic/gin"
)

func NewHealthzHandleFunc(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. check database
		defaultDBConfig := cfg.DatabaseMap["default"]
		dbConfigs := []config.Database{defaultDBConfig}
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

		// 4. return ok
		c.String(http.StatusOK, "ok")
	}
}
