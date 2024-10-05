package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/handler"
	"github.com/fengzhongzhu1621/xgo/ginx/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterVersion(cfg *config.Config, r *gin.RouterGroup) {

	// /api/v1/versions/*
	versionRouter := r.Group("versions")
	versionRouter.Use(middleware.Metrics())
	versionRouter.Use(middleware.AppLogger())
	{
		versionRouter.GET("/", handler.ListVersions(cfg))

		versionRouter.GET("/content/", handler.GetVersionContent(cfg))
	}
}
