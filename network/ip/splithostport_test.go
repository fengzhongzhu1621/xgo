package ip

import (
	"fmt"
	"net"
	"testing"
)

func TestSplitHostPort(t *testing.T) {
	addr := "example.com:8080"
	host, port, _ := net.SplitHostPort(addr)
	fmt.Println("Host:", host) // example.com
	fmt.Println("Port:", port) // 8080

	addr = "192.168.1.1:8080"
	host, port, _ = net.SplitHostPort(addr)
	fmt.Println("Host:", host) // 192.168.1.1
	fmt.Println("Port:", port) // 8080

	addr = "192.168.1.1"
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("Error:", err) // address 192.168.1.1: missing port in address
		return
	}
}
