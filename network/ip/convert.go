package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// IPtoUint32 IP从字符串转换为整型.
func IPtoUint32(ip string) uint32 {
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

// AddressToIPUint32 将IP:PORT格式的地址转换IP的整型.
func AddressToIPUint32(addr string) uint32 {
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	return IPtoUint32(ip)
}
