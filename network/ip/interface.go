package ip

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// GetIpByEthName 通过eth网卡名称获取机器ip.
func GetIpByEthName(name string) (string, error) {
	// 返回 interface 结构体对象的列表，包含了全部网卡信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
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

// GetRealIP 从请求头获取客户端 IP 地址
func GetRealIP(req *http.Request) string {
	xip := req.Header.Get("X-Real-IP")
	if xip == "" {
		xip = strings.Split(req.RemoteAddr, ":")[0]
	}
	return xip
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetAllIP 通过本机所有的IP.
func GetAllIP() []string {
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

// GetLocalConnectionIP 获得和远端服务器通信的IP地址.
func GetLocalConnectionIP(host string, port int) (string, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}

// Get preferred outbound ip of this machine
func GetPrivateIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ip = localAddr.IP.String()
	return
}

func GetPublicIP() (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	ip = string(body)
	return
}
