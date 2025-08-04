package auth

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestSecurityProvider(t *testing.T) {
	// 创建基本认证提供者
	basicAuth, err := NewSecurityProviderBasicAuth("user", "pass")
	if err != nil {
		fmt.Println("Basic Auth Error:", err)
		return
	}

	// 创建Bearer Token提供者
	bearerToken, err := NewSecurityProviderBearerToken("your_token_here")
	if err != nil {
		fmt.Println("Bearer Token Error:", err)
		return
	}

	// 创建API Key提供者，附加到Header
	apiKeyHeader, err := NewSecurityProviderApiKey("header", "X-API-Key", "your_api_key_here")
	if err != nil {
		fmt.Println("API Key Header Error:", err)
		return
	}

	// 创建API Key提供者，附加到Query参数
	apiKeyQuery, err := NewSecurityProviderApiKey("query", "api_key", "your_api_key_here")
	if err != nil {
		fmt.Println("API Key Query Error:", err)
		return
	}

	// 创建一个HTTP请求
	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		"https://api.example.com/data",
		nil,
	)
	if err != nil {
		fmt.Println("Request Creation Error:", err)
		return
	}

	// 使用各个提供者拦截请求并附加认证信息
	providers := []ISecurityProvider{basicAuth, bearerToken, apiKeyHeader, apiKeyQuery}
	for _, provider := range providers {
		if err := provider.Intercept(context.Background(), req); err != nil {
			fmt.Println("Intercept Error:", err)
			return
		}
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request Error:", err)
		return
	}
	defer resp.Body.Close()

	// 处理响应
	fmt.Println("Response Status:", resp.Status)
}
