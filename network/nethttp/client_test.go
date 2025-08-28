package nethttp

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	// 创建一个新的GET请求
	req, err := http.NewRequest("GET", "https://api.example.com/data", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求发送失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应内容失败:", err)
		return
	}

	// 输出响应内容
	fmt.Println("响应内容:", string(body))
}

// http.Client 有一个 Timeout 字段，用来控制整个请求的超时时间。
// 如果该时间内没有得到响应，请求就会自动取消，并返回超时错误。
// 该字段控制整个请求的超时时间，包括连接建立、发送请求、读取响应的总时间。如果超过设定时间，HTTP请求将超时。
func TestTimeout(t *testing.T) {
	// 设置请求超时时间为 5 秒
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 发送请求
	resp, err := client.Get("https://example.com")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}

func TestNewRequestWithContext(t *testing.T) {
	// 创建一个带有 5 秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return
	}

	// 创建一个 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}

// http.Transport 可以用来管理 HTTP 连接的底层实现，
// 通过设置 http.Transport 可以更细粒度地控制超时行为，比如设置连接超时、TLS握手超时、请求响应超时等。
func TestHttpTransport(*testing.T) {
	// 自定义 HTTP Transport
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second, // 设置连接建立的超时时间，避免因为网络连接问题导致的长时间等待。

		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,  // TLS 握手超时时间
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间

	}
	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // 总超时时间（包括连接、TLS握手、请求和响应）
	}

	// 发起 HTTP GET 请求
	resp, err := client.Get("https://example.com")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}
