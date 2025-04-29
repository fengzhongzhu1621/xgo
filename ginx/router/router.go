package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(cfg *config.Config, router *gin.Engine) {
	RegisterBasic(cfg, router)
	RegisterMetrics(cfg, router)
	RegisterPprof(cfg, router)
	RegisterSwagger(cfg, router)
	RegisterSchedule(cfg, router)
}

func RegisterRouterGroup(cfg *config.Config, rg *gin.RouterGroup) {
	RegisterVersion(cfg, rg)
}
