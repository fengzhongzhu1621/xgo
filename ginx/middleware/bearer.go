package middleware

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/fengzhongzhu1621/xgo/network/nethttp/auth/jwtx"
	"github.com/gin-gonic/gin"
)

// BearerAuthMiddleware 简单认证
func BearerAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 jwt token 的值
		bearerToken, err := jwtx.GetJwtTokenFromHeader(c)

		if err != nil {
			nethttp.UnauthorizedJSONResponse(
				c,
				err.Error())
			c.Abort()
			return
		}

		// 与配置文件中的 token 进行验证
		if bearerToken != "" && bearerToken != cfg.Auth.BearerToken {
			nethttp.UnauthorizedJSONResponse(
				c,
				"bearer token mismatch illegal")
			c.Abort()
			return
		}
		c.Next()
	}
}
