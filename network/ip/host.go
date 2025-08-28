package ip

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/fengzhongzhu1621/xgo/validator"
)

func getIP(addrPort string) (string, error) {
	idx := strings.LastIndex(addrPort, ":")
	return addrPort[:idx], nil
}

func getPort(addrPort string) (uint, error) {
	idx := strings.LastIndex(addrPort, ":")

	if len(addrPort[idx:]) < 2 {
		return 0, fmt.Errorf("the value of flag[AddrPort: %s] is wrong", addrPort)
	}
	port, err := strconv.ParseUint(addrPort[idx+1:], 10, 0)
	if err != nil {
		return uint(0), err
	}
	return uint(port), nil
}

func GetAddress(addrPort string) (string, error) {
	s := strings.TrimSpace(addrPort)
	if err := validator.CheckAddrPort(s); err != nil {
		return "", err
	}
	return getIP(s)
}

func GetPort(addrPort string) (uint, error) {
	s := strings.TrimSpace(addrPort)
	if err := validator.CheckAddrPort(s); err != nil {
		return uint(0), err
	}
	return getPort(s)
}

func GetDailAddress(addr string) (string, error) {
	uri, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	port := uri.Port()
	if uri.Port() == "" {
		port = "80"
	}
	return uri.Hostname() + ":" + port, err
}

// FilterHostIP filter out illegal IP addresses 过滤一个字符串切片，只保留其中有效的 IP 地址
func FilterHostIP(ipArr []string) []string {
	legalAddress := make([]string, 0)
	for _, address := range ipArr {
		// 接受 IPv4 和 IPv6 地址
		// 不会验证 IP 地址是否属于私有地址范围或保留地址
		// 如果输入是空字符串或其他非 IP 格式的字符串，会被过滤掉
		if ip := net.ParseIP(address); ip == nil {
			continue
		}
		legalAddress = append(legalAddress, address)
	}
	return legalAddress
}
