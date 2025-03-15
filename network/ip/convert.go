package ip

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"reflect"
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
	// 获得 ip 地址中的 ip 部分
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

// GetFirstIp 获得第一个 ip
func GetFirstIp(s string, sep string) string {
	if strings.Contains(s, ",") {
		return strings.Split(s, sep)[0]
	}

	return s
}

// ConvertHostIpv6Val convert host ipv6 value
func ConvertHostIpv6Val(items []string) ([]string, error) {
	var err error
	for idx, val := range items {
		items[idx], err = ConvertIPv6ToStandardFormat(val)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

// convertIpv6ToInt convert an ipv6 address to big int value
func convertIpv6ToInt(ipv6 string) (*big.Int, error) {
	ip := net.ParseIP(ipv6)
	if ip == nil {
		return nil, fmt.Errorf("invalid ipv6 address, data: %s", ipv6)
	}
	intVal := big.NewInt(0).SetBytes(ip)
	return intVal, nil
}

// convertIPv6ToFullAddr convert an ipv6 address to a full ipv6 address
func convertIPv6ToFullAddr(ipv6 string) (string, error) {
	if !strings.Contains(ipv6, ":") {
		return "", fmt.Errorf("address %s is not ipv6 address", ipv6)
	}

	intVal, err := convertIpv6ToInt(ipv6)
	if err != nil {
		return "", err
	}

	b255 := new(big.Int).SetBytes([]byte{255})
	buf := make([]byte, 2)
	part := make([]string, 8)
	pos := 0
	tmpInt := new(big.Int)
	var i uint
	for i = 0; i < 16; i += 2 {
		tmpInt.Rsh(intVal, 120-i*8).And(tmpInt, b255)
		bytes := tmpInt.Bytes()
		if len(bytes) > 0 {
			buf[0] = bytes[0]
		} else {
			buf[0] = 0
		}
		tmpInt.Rsh(intVal, 120-(i+1)*8).And(tmpInt, b255)
		bytes = tmpInt.Bytes()
		if len(bytes) > 0 {
			buf[1] = bytes[0]
		} else {
			buf[1] = 0
		}
		part[pos] = hex.EncodeToString(buf)
		pos++
	}

	return strings.Join(part, ":"), nil
}

// ConvertIPv6ToStandardFormat convert ipv6 address to standard format
// :: => 0000:0000:0000:0000:0000:0000:0000:0000
// ::127.0.0.1 => 0000:0000:0000:0000:0000:0000:127.0.0.1
func ConvertIPv6ToStandardFormat(address string) (string, error) {
	if ip := net.ParseIP(address); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	if !strings.Contains(address, ":") {
		return "", fmt.Errorf("address %s is not ipv6 address", address)
	}

	ipv6FullAddr, err := convertIPv6ToFullAddr(address)
	if err != nil {
		return "", err
	}

	addrs := strings.Split(address, ":")
	if !strings.Contains(addrs[len(addrs)-1], ".") {
		return ipv6FullAddr, nil
	}

	if ip := net.ParseIP(addrs[len(addrs)-1]); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	ipv6FullAddrs := strings.Split(ipv6FullAddr, ":")
	var result string
	for i := 0; i <= len(ipv6FullAddrs)-3; i++ {
		result += ipv6FullAddrs[i] + ":"
	}
	return result + addrs[len(addrs)-1], nil
}

// GetIPv4IfEmbeddedInIPv6 get ipv4 address if it is embedded in ipv6 address
// ::ffff:127.0.0.1 => 127.0.0.1, ::127.0.0.1 => 127.0.0.1
func GetIPv4IfEmbeddedInIPv6(address string) (string, error) {
	if ip := net.ParseIP(address); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	if !strings.Contains(address, ":") {
		return "", fmt.Errorf("address %s is not ipv6 address", address)
	}

	ipv6Addr, err := convertIPv6ToFullAddr(address)
	if err != nil {
		return "", err
	}
	ipv6Addrs := strings.Split(ipv6Addr, ":")
	for i := 0; i <= len(ipv6Addrs)-3; i++ {
		if i != len(ipv6Addrs)-3 && ipv6Addrs[i] != "0000" {
			return address, nil
		}

		if i == len(ipv6Addrs)-3 && ipv6Addrs[i] != "0000" && ipv6Addrs[i] != "ffff" {
			return address, nil
		}
	}

	addrs := strings.Split(address, ":")
	if !strings.Contains(addrs[len(addrs)-1], ".") {
		return address, nil
	}

	if ip := net.ParseIP(addrs[len(addrs)-1]); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	return addrs[len(addrs)-1], nil
}

// ConvertIpv6ToFullWord convert the ipv6 field into a complete format.
// for the converted scene at this time, there are only two types of value,
// string and slice, because the operators involved in mongo can only be
// one of the four cases "$eq", "$ne", "$in" and "$nin".
func ConvertIpv6ToFullWord(value interface{}) (interface{}, error) {
	var data interface{}
	switch reflect.ValueOf(value).Kind() {
	case reflect.String:
		ip, err := ConvertIPv6ToStandardFormat(value.(string))
		if err != nil {
			return nil, err
		}
		data = ip
	case reflect.Array, reflect.Slice:
		v := reflect.ValueOf(value)
		length := v.Len()
		if length == 0 {
			return value, nil
		}

		result := make([]interface{}, 0)
		// each element in the array or slice should be of the same basic type.
		for i := 0; i < length; i++ {
			item := v.Index(i).Interface()
			switch item.(type) {
			case string:
				v, err := ConvertIPv6ToStandardFormat(item.(string))
				if err != nil {
					return nil, err
				}
				result = append(result, v)
			default:
				return value, nil
			}
		}
		data = result
	default:
		return value, nil
	}

	return data, nil
}
