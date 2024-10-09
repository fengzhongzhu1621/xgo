package utils

import (
	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/gin-gonic/gin"
)

// GetRequestID 获得请求唯一 ID
func GetRequestID(c *gin.Context) string {
	return c.GetString(constant.RequestIDKey)
}

func SetRequestID(c *gin.Context, requestID string) {
	c.Set(constant.RequestIDKey, requestID)
}

func GetClientID(c *gin.Context) string {
	return c.GetString(constant.ClientIDKey)
}

func GetError(c *gin.Context) (interface{}, bool) {
	return c.Get(constant.ErrorIDKey)
}

func SetClientID(c *gin.Context, clientID string) {
	c.Set(constant.ClientIDKey, clientID)
}

func SetClientUsername(c *gin.Context, username string) {
	c.Set(constant.ClientUsernameKey, username)
}

func GetClientUsername(c *gin.Context) string {
	if name := c.GetString(constant.ClientUsernameKey); name != "" {
		return name
	} else {
		return constant.DefaultBackendOperator
	}
}

func SetBackendUser(c *gin.Context, user string) {
	c.Set(constant.BackendUserKey, user)
}

func GetBackendUser(c *gin.Context) string {
	if name := c.GetString(constant.BackendUserKey); name != "" {
		return name
	}
	return ""
}
