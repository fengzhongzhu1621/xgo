package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	log "github.com/sirupsen/logrus"
)

func CSRF(secret string) gin.HandlerFunc {
	log.Debug("Middleware: CSRF")

	// 使用adapter.Wrap方法将csrf.Protect包装起来，以便在Gin框架中使用
	return adapter.Wrap(
		// 使用密钥来创建 CSRF 保护
		// csrf.Secure(false): 设置CSRF保护为非安全模式，这意味着即使是在HTTPS连接中，CSRF令牌也可以通过非加密的HTTP连接传输。在生产环境中，你应该将其设置为true
		// csrf.Path("/")：设置CSRF保护的路径为根路径，这意味着所有路径都将受到保护。
		// csrf.CookieName 设置CSRF令牌的Cookie名称
		csrf.Protect(
			[]byte(secret),
			csrf.Secure(false),
			csrf.Path("/"),
			csrf.CookieName("xgo-csrf"),
		),
	)
}

func CSRFToken(domain string) gin.HandlerFunc {
	log.Debug("Middleware: CSRFToken")

	return func(c *gin.Context) {
		// 设置Cookie的SameSite属性为Lax，这意味着Cookie只会在从同一站点发送的（并且是GET类型的）请求中被发送。
		c.SetSameSite(http.SameSiteLaxMode)
		// csrf.Token(c.Request)生成的CSRF令牌
		// * 0：设置Cookie的最大生存时间为0，这意味着Cookie将在浏览器关闭时被删除。
		// * /：设置Cookie的作用域为整个站点。
		// * domain：设置Cookie的域名。
		// * false, false：设置Cookie为非安全和非HttpOnly模式。
		token := csrf.Token(c.Request)
		c.SetCookie("xgo-csrf-token", token, 0, "/", domain, false, false)

		c.Next()
	}
}
