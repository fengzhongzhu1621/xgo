package discovery

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/db/zookeeper/registerdiscover"
	"github.com/fengzhongzhu1621/xgo/ginx/backbone/configcenter"
	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
	log "github.com/sirupsen/logrus"
)

type IServiceManageInterface interface {
	// IsMaster 判断当前进程是否为master 进程， 服务注册节点的第一个节点
	// IsMaster() bool
	// Server(name string) Interface
	CoreService() Interface
}

type IDiscoveryInterface interface {
	IServiceManageInterface
}

type Interface interface {
	// GetServers 获取注册在zk上的所有服务节点
	GetServers() ([]string, error)
	// GetServersChan 最新的服务节点信息存放在该channel里，可被用来消费，以监听服务节点的变化
	GetServersChan() chan []string
}

// NewServiceDiscovery new a simple discovery module which can be used to get alive server address
func NewServiceDiscovery(client *registerdiscover.ZkClient, env string) (IDiscoveryInterface, error) {
	disc := registerdiscover.NewRegDiscoverEx(client)

	d := &discover{
		servers: make(map[string]*server),
	}

	curServiceName := server_option.GetIdentification()
	services := GetDiscoveryService()
	// 将当前服务也放到需要发现中
	services[curServiceName] = struct{}{}
	for component := range services {
		// 如果所有服务都按需发现服务。这个地方时不需要配置
		if component == server_option.MODULE_WEBSERVER && curServiceName != server_option.MODULE_WEBSERVER {
			continue
		}

		path := fmt.Sprintf("%s/%s", configcenter.SERV_BASEPATH, component)
		svr, err := newServerDiscover(disc, path, component, env)
		if err != nil {
			return nil, fmt.Errorf("discover %s failed, err: %v", component, err)
		}

		d.servers[component] = svr
	}

	return d, nil
}

type discover struct {
	servers map[string]*server
}

// Server 根据服务名获取服务再服务发现组件中的相关信息
func (d *discover) Server(name string) Interface {
	if svr, ok := d.servers[name]; ok {
		return svr
	}
	log.Infof("not found server. name: %s", name)

	return emptyServerInst
}

// IsMaster check whether current is master
func (d *discover) IsMaster() bool {
	return d.servers[server_option.GetIdentification()].IsMaster(server_info.GetServerInfo().UUID)
}

func (d *discover) CoreService() Interface {
	return d.servers[server_option.MODULE_CORESERVICE]
}
