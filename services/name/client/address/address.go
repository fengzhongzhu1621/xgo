package address

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type INameServiceAddress interface {
	GetAddress() string
	ParseAddress(addr string) error
}

/**
 * 名字服务的地址
 * 名字服务：以SID（由模块ID和命令字ID组成）为关键字，通过SID取得真正的IP和端口地址，使得IP和端口配置对调用者透明
 */
type NameServiceAddress struct {
	Type    string   // 地址类型
	Address string   // 具体的地址
	IpList  []string // 多个IP地址
	ModId   int
	CmdId   int
}

func getRandomIp(IpList []string) string {
	return IpList[rand.Int()%len(IpList)]
}

func (addr *NameServiceAddress) GetAddress() string {
	if addr.Type == "http" {
		if strings.HasPrefix(addr.Address, "http://") {
			return addr.Address
		} else {
			return fmt.Sprintf("http://%s", addr.Address)
		}
	} else if addr.Type == "ip" {
		// 随机从IP列表中获取一个地址
		return fmt.Sprintf("ip://%s", getRandomIp(addr.IpList))
	} else if addr.Type == "l5" {
		if strings.HasPrefix(addr.Address, "l5://") {
			return addr.Address
		} else {
			return fmt.Sprintf("l5://%s", addr.Address)
		}
	}
	return ""
}

func (addr *NameServiceAddress) ParseAddress(address string) error {
	if strings.HasPrefix(address, "http://") {
		addr.Type = "http"
		addr.Address = address
		return nil
	} else if strings.HasPrefix(address, "ip://") {
		addr.Type = "ip"
		addr.IpList = strings.Split(address[5:], ",")
		addr.Address = address
		return nil
	} else if strings.HasPrefix(address, "l5://") {
		addr.Type = "l5"
		l5Ids := strings.Split(address[5:], ":")
		addr.ModId, _ = strconv.Atoi(l5Ids[0])
		addr.CmdId, _ = strconv.Atoi(l5Ids[1])
		addr.Address = address
		return nil
	}
	return errors.New("invalid address")
}
