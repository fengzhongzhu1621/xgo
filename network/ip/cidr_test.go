package ip

import (
	"fmt"
	"net"
	"testing"
)

func TestParseCIDR(t *testing.T) {
	cidr := "192.168.1.0/24"

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("解析 CIDR 失败:", err)
		return
	}

	fmt.Printf("解析后的网络地址: %s\n", ipNet.IP) // 192.168.1.0
	fmt.Printf("子网掩码: %s\n", ipNet.Mask)   // ffffff00
}
