package utils

import (
	"os"

	"github.com/gin-gonic/gin"
)

func SetGinMode() {
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		gin.SetMode(gin.ReleaseMode)
		return
	}
	gin.SetMode(mode)
}
