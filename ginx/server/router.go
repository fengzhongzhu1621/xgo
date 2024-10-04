package server

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DisableConsoleColor()

	router := gin.New()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(middleware.Recovery(cfg.Sentry.Enable))
	router.Use(middleware.RequestID())

	return router
}
