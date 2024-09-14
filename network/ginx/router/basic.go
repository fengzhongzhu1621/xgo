package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/network/ginx/handler"
	"github.com/gin-gonic/gin"
)

// RegisterBasic 注册健康检查和版本基础的路由
func RegisterBasic(cfg *config.Config, router *gin.Engine) {
	router.GET("/ping", handler.Ping)
	router.GET("/healthz", handler.NewHealthzHandleFunc(cfg))
	router.GET("/version", handler.Version)
}
