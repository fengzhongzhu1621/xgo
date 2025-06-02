package dns

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

// FetchSubdomains 获取指定域名的子域名列表
func FetchSubdomains(domain string) []string {
	cmd := exec.Command("subfinder", "-d", domain, "-silent")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Subfinder error: %v", err)
	}
	lines := strings.Split(string(out), "\n")
	var results []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			results = append(results, line)
		}
	}
	return results
}

// ResolveDomain 接受一个域名字符串，返回该域名的第一个IP地址和对应的CNAME记录。
func ResolveDomain(domain string) (string, string) {
	// 使用net.LookupIP查询域名的所有IP地址
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("查询IP地址失败: %v\n", err)
		return "", ""
	}

	// 使用net.LookupCNAME查询域名的CNAME记录
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		fmt.Printf("查询CNAME记录失败: %v\n", err)
		return "", ""
	}

	// 初始化一个变量用于存储第一个IP地址的字符串形式
	var ipStr string

	// 检查是否成功获取到IP地址
	if len(ips) > 0 {
		// 将第一个IP地址转换为字符串并赋值给ipStr
		ipStr = ips[0].String()
	}

	// 返回第一个IP地址和CNAME记录
	return ipStr, cname
}
