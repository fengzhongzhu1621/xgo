package header

import (
	"strings"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/gin-gonic/gin"
)

// GetLanguageFromHeader 从请求头获取语言
func GetLanguageFromHeader(c *gin.Context) (string, error) {
	header := c.Request.Header.Get(constant.HTTPHeadLanguage)
	if len(header) == 0 {
		return "", xgo.JwtTokenNoneErr
	}
	strs := strings.Split(header, " ")

	return strs[0], nil
}

// GetEnvFromHeader 从请求头获取 env 的值，如果找不到则从 Get 请求参数中获取
func GetEnvFromHeader(c *gin.Context) string {
	env := c.Request.Header.Get("env")
	if env == "" {
		env = c.Query("env")
	}

	return env
}

// GetUsernameFromHeader 从 header 中获取用户名
func GetUsernameFromHeader(c *gin.Context) string {
	username := c.Request.Header.Get("bk_username")
	if username == "" {
		username = c.Query("bk_username")
	}

	return username
}

// GetTokenFromHeader 从 header 中获取 token
func GetTokenFromHeader(c *gin.Context) string {
	token := c.Request.Header.Get("token")
	if token == "" {
		token = c.Query("token")
	}

	return token
}
