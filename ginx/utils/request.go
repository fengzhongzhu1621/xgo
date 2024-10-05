package utils

import (
	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/gin-gonic/gin"
)

// GetRequestID 获得请求唯一 ID
func GetRequestID(c *gin.Context) string {
	return c.GetString(constant.RequestIDKey)
}

// SetRequestID ...
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(constant.RequestIDKey, requestID)
}

func GetClientID(c *gin.Context) string {
	return c.GetString(constant.ClientIDKey)
}

func GetError(c *gin.Context) (interface{}, bool) {
	return c.Get(constant.ErrorIDKey)
}
