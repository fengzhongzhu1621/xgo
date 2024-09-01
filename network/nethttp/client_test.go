package nethttp

import (
	"fmt"
	"io"
	"net/http"
	"testing"
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
