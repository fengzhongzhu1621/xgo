package dns

import (
	"fmt"
	"testing"
)

func TestFetchSubdomains(t *testing.T) {
	domain := "www.example.com"

	subDomains := FetchSubdomains(domain)
	fmt.Printf("域名: %v\n", subDomains)
}

func TestResolveDomain(t *testing.T) {
	// 示例域名
	domain := "www.example.com"

	// 调用ResolveDomain函数获取IP和CNAME
	ip, cname := ResolveDomain(domain)

	// 打印结果
	fmt.Printf("域名: %s\n", domain)
	fmt.Printf("第一个IP地址: %s\n", ip)
	fmt.Printf("CNAME记录: %s\n", cname)
}
