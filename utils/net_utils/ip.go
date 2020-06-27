package net_utils

import (
	"fmt"
	"net"
	"strings"
)

/**
 * 通过eth网卡名称获取机器ip
 */
func GetIpv4ByName(name string) (string, error) {
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
						if ipnet.IP.To4() != nil {
							return ipnet.IP.String(), nil
						}
					}
				}
			}
		}
	}
	return "", nil
}

/**
 * 通过本机所有的IP
 */
func GetAllIpv4() []string {
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
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
		}
	}
	return ips
}

/**
 * 根据主机和端口获取IP
 */
func GetLocalIp(host string, port int) (string, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}
