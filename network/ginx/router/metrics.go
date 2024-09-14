package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetrics 注册 metrics 指标路由
func RegisterMetrics(cfg *config.Config, router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
