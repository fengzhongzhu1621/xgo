package validator

import (
	"fmt"
	"net"
	"testing"

	"github.com/duke-git/lancet/v2/netutil"
	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// TestIsDns 检查字符串是否是有效的域名系统（DNS）。
// func IsDns(dns string) bool
func TestIsDns(t *testing.T) {
	result1 := validator.IsDns("abc.com")
	result2 := validator.IsDns("a.b.com")
	result3 := validator.IsDns("http://abc.com")
	result4 := validator.IsDns("http://abc.com?a=1")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, true, result2, "result2")
	assert.Equal(t, false, result3, "result3")
	assert.Equal(t, false, result4, "result4")
}

// TestIsUrl 检查字符串是否是有效的URL
// func IsUrl(str string) bool
func TestIsUrl(t *testing.T) {
	result1 := validator.IsUrl("abc.com")
	result2 := validator.IsUrl("http://abc.com")
	result3 := validator.IsUrl("abc")
	result4 := validator.IsUrl("http://abc.com?a=1")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, true, result2, "result2")
	assert.Equal(t, false, result3, "result3")
	assert.Equal(t, true, result4, "result4")
}

// TestLancetIsIp 检查字符串是否为IP地址。
// func IsIp(ipstr string) bool
func TestLancetIsIp(t *testing.T) {
	result1 := validator.IsIp("127.0.0.1")
	result2 := validator.IsIp("::0:0:0:0:0:0:1")
	result3 := validator.IsIp("127.0.0")
	result4 := validator.IsIp("::0:0:0:0:")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, true, result2, "result2")
	assert.Equal(t, false, result3, "result3")
	assert.Equal(t, false, result4, "result4")
}

// TestIsIpV4 检查字符串是否为IPv4地址。
// func IsIpV4(ipstr string) bool
func TestIsIpV4(t *testing.T) {
	result1 := validator.IsIpV4("127.0.0.1")
	result2 := validator.IsIpV4("::0:0:0:0:0:0:1")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, false, result2, "result2")
}

// TestIsIpV6 检查字符串是否为IPv6地址。
// func TestIsIpV6(ipstr string) bool
func TestIsIpV6(t *testing.T) {
	result1 := validator.IsIpV6("127.0.0.1")
	result2 := validator.IsIpV6("::0:0:0:0:0:0:1")

	assert.Equal(t, false, result1, "result1")
	assert.Equal(t, true, result2, "result2")
}

// TestTestIsPublicIP Checks if an ip is public or not.
// func IsPublicIP(IP net.IP) bool
func TestIsPublicIP(t *testing.T) {
	ip1 := netutil.IsPublicIP(net.ParseIP("127.0.0.1"))
	ip2 := netutil.IsPublicIP(net.ParseIP("192.168.0.1"))
	ip3 := netutil.IsPublicIP(net.ParseIP("36.112.24.10"))

	assert.Equal(t, false, ip1, "ip1")
	assert.Equal(t, false, ip2, "ip2")
	assert.Equal(t, true, ip3, "ip3")
}

// TestIsInternalIP Checks if an ip is intranet or not.
// func IsInternalIP(IP net.IP) bool
func TestIsInternalIP(t *testing.T) {
	ip1 := netutil.IsPublicIP(net.ParseIP("127.0.0.1"))
	ip2 := netutil.IsPublicIP(net.ParseIP("192.168.0.1"))
	ip3 := netutil.IsPublicIP(net.ParseIP("36.112.24.10"))

	assert.Equal(t, true, ip1, "ip1")
	assert.Equal(t, true, ip2, "ip2")
	assert.Equal(t, false, ip3, "ip3")
}

// TestIsPingConnected checks if can ping the specified host or not.
// func IsPingConnected(host string) bool
func TestIsPingConnected(t *testing.T) {
	result1 := netutil.IsPingConnected("www.baidu.com")
	result2 := netutil.IsPingConnected("www.!@#&&&.com")

	// panic: test timed out after 30s
	fmt.Println(result1)
	fmt.Println(result2)
}

// TestIsTelnetConnected Checks if can telnet the specified host or not.
// func IsTelnetConnected(host string, port string) bool
func TestIsTelnetConnected(t *testing.T) {
	result1 := netutil.IsTelnetConnected("www.baidu.com", "80")
	result2 := netutil.IsTelnetConnected("www.baidu.com", "123")

	fmt.Println(result1) // true
	fmt.Println(result2) // true
}
