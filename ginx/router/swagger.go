package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterSwagger 注册文档
func RegisterSwagger(cfg *config.Config, router *gin.Engine) {
	if cfg.EnableSwagger {
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
}
