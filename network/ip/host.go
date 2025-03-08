package ip

import (
	"fmt"
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
