package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/fengzhongzhu1621/xgo/version"
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
		"version":   version.Version,
		"commit":    version.Commit,
		"buildTime": version.BuildTime,
		"goVersion": version.GoVersion,
		"env":       runEnv,
		"timestamp": now.Unix(),
		"date":      now,
	})
}
