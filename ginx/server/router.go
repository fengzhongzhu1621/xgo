package server

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/middleware"
	"github.com/fengzhongzhu1621/xgo/ginx/router"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DisableConsoleColor()

	rootRouter := gin.New()

	// 添加中间件
	rootRouter.Use(gin.Logger())
	rootRouter.Use(middleware.Recovery(cfg.Sentry.Enable))
	rootRouter.Use(middleware.RequestID())

	// 注册默认路由
	router.RegisterRouter(cfg, rootRouter)

	// 注册自定义路由组
	router.RegisterRouterGroup(cfg, rootRouter.Group("/api/v1"))

	return rootRouter
}
