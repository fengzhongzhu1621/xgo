package net_utils

import (
	"bytes"
	"encoding/binary"
	"net"
	"strings"
)

// 通过eth网卡名称获取机器ip
func GetIpByEthName(name string) (string, error) {
	// 返回 interface 结构体对象的列表，包含了全部网卡信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", nil
	}

	for _, v := range interfaces {
		if strings.HasPrefix(v.Name, name) {
			// 通过判断net.FlagUp标志进行确认，排除掉无用的网卡
			if (v.Flags & net.FlagUp) != 0 {
				return "", nil
			}
			// 返回一个网卡上全部的IP列表
			if addrs, err := v.Addrs(); err == nil {
				for _, address := range addrs {
					if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}
	return "", nil
}

// 通过本机所有的IP
func GetAllIp() []string {
	// 返回 interface 结构体对象的列表，包含了全部网卡信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	ips := make([]string, 0)

	for _, v := range interfaces {
		// 通过判断net.FlagUp标志进行确认，排除掉无用的网卡
		if (v.Flags & net.FlagUp) != 0 {
			return nil
		}
		// 返回一个网卡上全部的IP列表
		if addrs, err := v.Addrs(); err == nil {
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}
	return ips
}

// 获得和远端服务器通信的IP地址
func GetLocalConnectionIp(host string, port int) (string, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}

// IP从字符串转换为整型
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

// IP从整型转换为字符串
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

/**
 * 整型大小端互转
 *   3    2     1    0
 *  XX   @@     ##   $$
 *         ||
 *         \/
 *   3    2     1    0
 *  $$   ##     @@   XX
 */
func ConvertEndian(num uint32) uint32 {
	return ((num >> 24) & 0xff) |
		((num << 8) & 0xff0000) |
		((num >> 8) & 0xff00) |
		((num << 24) & 0xff000000)
}

// 将IP:PORT格式的地址转换IP的整型
func AddressToIpUint32(addr string) uint32 {
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	return IPtoUint32(ip)
}
