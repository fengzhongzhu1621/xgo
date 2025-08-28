package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// ProxyHttp 一个基于 Gin 框架的 HTTP 反向代理功能，主要作用是将接收到的 HTTP 请求转发到指定的目标地址 addr
// addr 目标服务器的地址（如 "http://example.com" 或 "https://api.example.com"）。
func ProxyHttp(c *gin.Context, tlsConf *tls.Config, addr string) {
	// 使用 url.Parse 解析目标地址 addr，如果解析失败，直接返回错误信息给客户端。
	u, err := url.Parse(addr)
	if err == nil {
		// 创建一个反向代理对象，所有请求会被转发到解析后的目标地址 u（如 http://example.com）。
		proxy := httputil.NewSingleHostReverseProxy(u)
		if tlsConf != nil {
			// 自定义 HTTP 传输层配置，用于控制代理的底层连接行为
			proxy.Transport = &http.Transport{
				// 使用环境变量中的 HTTP 代理设置
				Proxy: http.ProxyFromEnvironment,
				// 配置 TCP 连接参数（超时和保活）
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				// 强制尝试 HTTP/2
				ForceAttemptHTTP2: true,
				// 最大空闲连接数（100）
				MaxIdleConns: 100,
				// 空闲连接超时时间（90秒）
				IdleConnTimeout: 90 * time.Second,
				// TLS 握手超时时间（10秒）
				TLSHandshakeTimeout: 10 * time.Second,
				// 用于控制客户端在发送 HTTP 请求时，等待服务器响应 100 Continue 状态码的超时时间。
				// 定义了客户端在发送请求头后，等待服务器 100 Continue 响应的最长时间。
				// 如果超时时间内未收到响应，客户端会直接发送请求体（或中止请求，取决于实现）。
				// 默认值为 1 秒（在 Go 的 http.Transport 中）。
				// 如果请求未设置 Expect: 100-continue 头部，此参数无效。
				// 对于小请求体，通常不需要启用 100-continue，因为它会增加一次 RTT（往返时间）。
				// 通过合理设置 ExpectContinueTimeout，可以平衡请求的可靠性和性能。
				//
				// HTTP 的 Expect: 100-continue 机制
				// 当客户端发送一个带有 Expect: 100-continue 头部的请求时，客户端会先发送请求头（不包括请求体），并等待服务器返回 100 Continue 响应。
				// 如果服务器同意接收请求体，客户端才会继续发送请求体；否则，客户端可以中止发送（例如，当服务器返回 417 Expectation Failed 时）。
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig:       tlsConf, // 指定 TLS 配置（如证书、密钥等），用于 HTTPS 请求
			}
		}

		// 将 Gin 的请求（c.Request）和响应写入器（c.Writer）交给反向代理处理，完成请求的转发和响应的回传
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.Writer.Write([]byte(err.Error()))
	}
}
