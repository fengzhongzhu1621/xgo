package webservice

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/db/redis/session"
	"github.com/fengzhongzhu1621/xgo/ginx/backbone"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Service struct {
	*server_option.ServerOption
	Engine   *backbone.Engine
	CacheCli redis.Client
	Config   *config.WebServerConfig
	Session  session.RedisStore
}

func (s *Service) WebService() *gin.Engine {
	utils.SetGinMode()
	ws := gin.New()
	ws.Use(gin.Logger())

	return ws
}
