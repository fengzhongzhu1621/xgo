package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	_ "github.com/fengzhongzhu1621/xgo/ginx/docs" // 引入 docs 包，以使 swag 自动生成文档
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterSwagger 注册文档
func RegisterSwagger(cfg *config.Config, router *gin.Engine) {
	if cfg.EnableSwagger {
		// 设置 Swagger 文档路由
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
