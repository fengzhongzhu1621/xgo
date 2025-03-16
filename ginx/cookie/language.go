package cookie

import (
	"github.com/fengzhongzhu1621/xgo/network/constant"
	"github.com/gin-gonic/gin"
)

// GetLanguageByHTTPRequest 从客户端获得语言
func GetLanguageFromCookie(c *gin.Context) string {
	cookieLanguage, err := c.Cookie(constant.HTTPCookieLanguage)
	if err == nil && cookieLanguage != "" {
		return cookieLanguage
	}

	return string(constant.Chinese)
}
