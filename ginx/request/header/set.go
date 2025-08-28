package header

import (
	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/fengzhongzhu1621/xgo/ginx/cookie"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetProxyHeader(c *gin.Context) {
	// http request header add user
	session := sessions.Default(c)
	userName, _ := session.Get(constant.WEBSessionUinKey).(string)
	ownerID, _ := session.Get(constant.WEBSessionOwnerUinKey).(string)

	// 删除 Accept-Encoding 避免返回值被压缩
	c.Request.Header.Del("Accept-Encoding")

	nethttp.AddUser(c.Request.Header, userName)
	nethttp.AddLanguage(c.Request.Header, cookie.GetLanguageFromCookie(c))
	nethttp.AddSupplierAccount(c.Request.Header, ownerID)
}
