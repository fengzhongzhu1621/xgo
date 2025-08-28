package router

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/handler"
	"github.com/gin-gonic/gin"
)

// RegisterSchedule 注册cron调度任务路由
func RegisterSchedule(cfg *config.Config, router *gin.Engine) {
	router.POST("/task", handler.AddTask)
	router.POST("/task/:id/run", handler.RunNow)
	router.POST("/task/:id/pause", handler.PauseTask)
	router.POST("/task/:id/resume", handler.ResumeTask)
}
