package auth

import (
	"strings"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/gin-gonic/gin"
)

// GetJwtTokenFromHeader 从请求头获取 jwt 的值
func GetJwtTokenFromHeader(c *gin.Context) (string, error) {
	// 找不到认证头部
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return "", xgo.JwtTokenNoneErr
	}
	// 获得 Bearer 的值
	strs := strings.SplitN(header, " ", 2)
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", xgo.JwtTokenInvalidErr
	}

	return strs[1], nil
}
