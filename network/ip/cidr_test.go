package ip

import (
	"fmt"
	"net"
	"testing"
)

func TestParseCIDR(t *testing.T) {
	cidr := "192.168.1.0/24"

	// *IPNet：一个指向 IPNet 结构体的指针，表示解析后的网络地址和子网掩码。
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("解析 CIDR 失败:", err)
		return
	}
	fmt.Printf("解析后的网络地址: %s\n", ipNet.IP) // 192.168.1.0
	fmt.Printf("子网掩码: %s\n", ipNet.Mask)   // ffffff00

	cidr = "192.168.1.10/24"
	_, ipNet, _ = net.ParseCIDR(cidr)
	fmt.Printf("解析后的网络地址: %s\n", ipNet.IP) // 192.168.1.0
	fmt.Printf("子网掩码: %s\n", ipNet.Mask)   // ffffff00
}
