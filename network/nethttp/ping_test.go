package nethttp

import (
	"fmt"
	"testing"
)

func TestIsHTTPAlive(t *testing.T) {
	// 示例域名
	domains := []string{
		"www.google.com",
		"www.nonexistentwebsite123.com",
		"example.com",
	}

	// 遍历每个域名并检查其HTTP状态
	for _, domain := range domains {
		if IsHTTPAlive(domain) {
			fmt.Printf("域名 %s 在线\n", domain)
		} else {
			fmt.Printf("域名 %s 不在线或无法访问\n", domain)
		}
	}
}
