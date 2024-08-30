package ip

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Inet4ToInt64 net.Ip 类型转换为 int64 类型
func Inet4ToInt64(ip net.IP) int64 {
	ipv4Int := big.NewInt(0)
	ipv4Int.SetBytes(ip.To4())
	return ipv4Int.Int64()
}

// Ipv4ToInt64 ipv4 字符串转换为 int64 类型
func Ipv4ToInt64(ip string) int64 {
	ipt := net.ParseIP(ip)
	if ipt == nil {
		return 0
	}
	return Inet4ToInt64(ipt)
}

// InetToHexStr net.Ip 类型转换为十六进制字符串
func InetToHexStr(ip net.IP) string {
	ipv4 := false
	// 判断是否是 ipv4
	if ip.To4() != nil {
		ipv4 = true
	}

	ipInt := big.NewInt(0)
	if ipv4 {
		ipInt.SetBytes(ip.To4())
		ipHex := hex.EncodeToString(ipInt.Bytes())
		return ipHex
	}

	ipInt.SetBytes(ip.To16())
	ipHex := hex.EncodeToString(ipInt.Bytes())
	return ipHex
}

// IpToUint32 IP从字符串转换为整型.
func IpToUint32(ip string) uint32 {
	// 检查IP地址格式是否有效
	// 解析为IP地址，并返回该地址。如果s不是合法的IP地址文本表示，ParseIP会返回nil
	ips := net.ParseIP(ip)
	if ips == nil {
		return 0
	}
	if len(ips) == net.IPv6len {
		return binary.BigEndian.Uint32(ips[12:16])
	} else if len(ips) == net.IPv4len {
		return binary.BigEndian.Uint32(ips)
	}
	return 0
}

// Uint32toIP IP从整型转换为字符串.
func Uint32toIP(ip uint32) string {
	// 整型转换为字节序列
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, ip)
	if err != nil {
		return ""
	}
	b := buf.Bytes()
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

/*
 ConvertEndian 整型大小端互转.
    3    2     1    0
   XX   @@     ##   $$
          ||
          \/
    3    2     1    0
   $$   ##     @@   XX
*/
func ConvertEndian(num uint32) uint32 {
	return ((num >> 24) & 0xff) |
		((num << 8) & 0xff0000) |
		((num >> 8) & 0xff00) |
		((num << 24) & 0xff000000)
}

// AddressToIpUint32 将IP:PORT格式的地址转换IP的整型.
func AddressToIpUint32(addr string) uint32 {
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	return IpToUint32(ip)
}

// CheckIpCidr 判断 ip 是否在 ip 段列表中
func CheckIpInCidrList(srcIp string, cidrList []string) bool {
	// 活的 ip 地址中的 ip 部分
	var ip string
	if strings.Contains(srcIp, ":") {
		host, _, err := net.SplitHostPort(srcIp)
		if err != nil {
			return false
		}
		ip = host
	} else {
		ip = srcIp
	}
	// 将 ip 字符串转换为 net.IP
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return false
	}

	// 判断 ip 是否在 ip 段中
	for _, v := range cidrList {
		// 获得 ip 地址范围
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			continue
		}
		if ipNet.Contains(netIp) {
			return true
		}
	}

	return false
}

// ToCanonicalIP replaces ":0:0" in IPv6 addresses with "::"
// ToCanonicalIP("192.168.0.1") -> "192.168.0.1"
// ToCanonicalIP("100:200:0:0:0:0:0:1") -> "100:200::1".
func ToCanonicalIP(host string) string {
	val := net.ParseIP(host)
	if val == nil {
		return host
	}
	return val.String()
}

// ExpandIP expands IPv6 addresses "::" to ":0:0..."
func ExpandIP(host string) string {
	if !strings.Contains(host, "::") {
		return host
	}
	expected := 7
	existing := strings.Count(host, ":") - 1
	return strings.Replace(host, "::", strings.Repeat(":0", expected-existing)+":", 1)
}
