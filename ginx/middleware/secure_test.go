package middleware

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/unrolled/secure"
)

// secure 中间件增加了安全检查层，可能会引入轻微的延迟。确保您的服务器经过优化，能够处理额外的处理需求

func TestSecureHttpMiddleware(t *testing.T) {
	port := os.Getenv("PORT")

	secureMiddleware := secure.New(secure.Options{
		// 指定允许访问的域名，防止未授权的主机访问
		AllowedHosts: []string{"example.com", "ssl.example.com"},
		// 强制将 HTTP 重定向到 HTTPS
		// 通过 SSLRedirect: true 和 SSLHost: "ssl.example.com"，所有 HTTP 请求将被重定向到 HTTPS。
		SSLRedirect: true,
		// 指定 SSL/TLS 证书绑定的域名
		// 确保 ssl.example.com 的 SSL 证书已正确配置。可以使用自签名证书进行开发和测试，但在生产环境中应使用受信任的证书颁发机构（CA）签发的证书。
		SSLHost: "ssl.example.com",
		// 设置 HSTS，要求浏览器在指定时间内只通过 HTTPS 访问站点（315360000 秒 = 1 年）
		STSSeconds: 315360000,
		// HSTS 策略同样应用于子域名
		STSIncludeSubdomains: true,
		// 禁止页面在 frame 中展示，防止点击劫持
		FrameDeny: true,
		// 禁止浏览器猜测内容类型，防止 MIME 类型混淆攻击
		ContentTypeNosniff: true,
		// 启用浏览器内置的 XSS 防护
		BrowserXssFilter: true,
		// 设置内容安全策略 (CSP)，限制资源加载来源
		ContentSecurityPolicy: "default-src 'self'",
	})

	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	handler := secureMiddleware.Handler(app)
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func TestSecureInGin(t *testing.T) {
	router := gin.Default()

	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     "localhost:8080",
	})

	router.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)
			if err != nil {
				c.Abort()
				return
			}
			c.Next()
		}
	}())

	router.Run(":8080")
}
