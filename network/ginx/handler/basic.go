package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Version(c *gin.Context) {
	runEnv := os.Getenv("RUN_ENV")
	now := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"version":   xgo.Version,
		"commit":    xgo.Commit,
		"buildTime": xgo.BuildTime,
		"goVersion": xgo.GoVersion,
		"env":       runEnv,
		"timestamp": now.Unix(),
		"date":      now,
	})
}
