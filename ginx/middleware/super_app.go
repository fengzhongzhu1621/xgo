package middleware

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/gin-gonic/gin"
)

func SuperClientAppMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		appCode := utils.GetClientID(c)
		if !config.SuperAppCodeSet.Has(appCode) {
			nethttp.UnauthorizedJSONResponse(c, "super client app code wrong")
			c.Abort()
			return
		}

		c.Next()
	}
}