package middleware

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/fengzhongzhu1621/xgo/network/nethttp/auth/jwtx"
	"github.com/gin-gonic/gin"
)

// BackendAuthMiddleware 后端服务通信的鉴权中间件
func BackendAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 jwt 的值
		jwt_token, err := jwtx.GetJwtTokenFromHeader(c)

		if err != nil {
			utils.UnauthorizedJSONResponse(
				c,
				err.Error())
			c.Abort()
			return
		}

		// 解析 jwt token
		var (
			option jwtx.CustomJwtClaimsOption
			claims *jwtx.CustomJwtClaims
		)
		option.Cfg = cfg
		option.HS256Key = cfg.Auth.JwtToken
		claims, _ = jwtx.NewCustomJwtClaims(&option)
		jwtclaims, err := claims.ParseHS256JwtToken(jwt_token)

		if err != nil {
			utils.UnauthorizedJSONResponse(
				c,
				err.Error())
			c.Abort()
			return
		}

		// 解析 jwt token 获取后台服务通信用的用户名
		utils.SetBackendUser(c, jwtclaims.GetOperator())
		c.Next()
	}
}
