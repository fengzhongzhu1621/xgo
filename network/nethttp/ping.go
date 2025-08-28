package nethttp

import (
	"net/http"
	"time"
)

// IsHTTPAlive 检查指定的域名是否可以通过HTTP访问
func IsHTTPAlive(domain string) bool {
	// 创建一个带有超时设置的HTTP客户端
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时时间为5秒
	}

	// 发送GET请求到指定的域名
	resp, err := client.Get("http://" + domain)
	if err != nil {
		// 如果请求过程中发生错误，返回false
		return false
	}
	// 确保响应体在函数结束时被关闭
	defer resp.Body.Close()

	// 如果响应状态码小于500，认为服务在线
	return resp.StatusCode < 500
}
