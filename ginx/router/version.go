package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/handler"
	"github.com/gin-gonic/gin"
)

func RegisterVersion(cfg *config.Config, r *gin.RouterGroup) {

	// /api/v1/web/versions/*
	versionRouter := r.Group("versions")
	{
		versionRouter.GET("/", handler.ListVersions(cfg))

		versionRouter.GET("/content/", handler.GetVersionContent(cfg))
	}
}