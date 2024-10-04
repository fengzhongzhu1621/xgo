package router

import (
	"net/http/pprof"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/gin-gonic/gin"
)

// RegisterPprof 注册 pprof 路由
func RegisterPprof(cfg *config.Config, router *gin.Engine) {
	pprofRouter := router.Group("/debug/pprof")
	if !cfg.Debug {
		pprofRouter.Use(gin.BasicAuth(cfg.PProf.Account))
	}

	pprofRouter.GET("/", gin.WrapF(pprof.Index))
	pprofRouter.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	pprofRouter.GET("/profile", gin.WrapF(pprof.Profile))
	pprofRouter.POST("/symbol", gin.WrapF(pprof.Symbol))
	pprofRouter.GET("/symbol", gin.WrapF(pprof.Symbol))
	pprofRouter.GET("/trace", gin.WrapF(pprof.Trace))
	pprofRouter.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
	pprofRouter.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	pprofRouter.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	pprofRouter.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	pprofRouter.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	pprofRouter.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))

}
