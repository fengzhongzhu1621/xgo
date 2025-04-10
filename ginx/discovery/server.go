package discovery

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fengzhongzhu1621/xgo/db/zookeeper/registerdiscover"
	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
	log "github.com/sirupsen/logrus"
)

type server struct {
	sync.RWMutex
	next int
	// server's name
	name         string
	path         string
	environment  string
	servers      []*server_info.ServerInfo
	master       *server_info.ServerInfo
	discoverChan <-chan *registerdiscover.DiscoverEvent
	serversChan  chan []string
}

func newServerDiscover(disc *registerdiscover.RegDiscover, path, name, env string) (*server, error) {
	discoverChan, eventErr := disc.DiscoverService(path)
	if nil != eventErr {
		return nil, eventErr
	}

	svr := &server{
		path:         path,
		name:         name,
		environment:  env,
		servers:      make([]*server_info.ServerInfo, 0),
		discoverChan: discoverChan,
		serversChan:  make(chan []string, 1),
	}

	svr.run()
	return svr, nil
}

// GetServersChan 获取zk上最新的服务节点信息channel
func (s *server) GetServersChan() chan []string {
	return s.serversChan
}

// getInstances 获取所有注册服务节点的ip:port
func (s *server) getInstances() []string {
	addrArr := []string{}
	s.RLock()
	defer s.RUnlock()
	for _, info := range s.servers {
		addrArr = append(addrArr, info.Instance())
	}
	return addrArr
}

func (s *server) GetServers() ([]string, error) {
	if s == nil {
		return []string{}, nil
	}

	s.Lock()
	num := len(s.servers)
	if num == 0 {
		s.Unlock()
		return []string{}, fmt.Errorf("oops, there is no %s can be used", s.name)
	}

	var infos []*server_info.ServerInfo
	if s.next < num-1 {
		s.next = s.next + 1
		infos = append(s.servers[s.next-1:], s.servers[:s.next-1]...)
	} else {
		s.next = 0
		infos = append(s.servers[num-1:], s.servers[:num-1]...)
	}
	s.Unlock()

	servers := make([]string, 0)
	for _, server := range infos {
		servers = append(servers, server.RegisterAddress())
	}

	return servers, nil
}

// IsMaster 判断当前进程是否为master 进程， 服务注册节点的第一个节点
// 注册地址不能作为区分标识，因为不同的机器可能用一样的域名作为注册地址，所以用uuid区分
func (s *server) IsMaster(UUID string) bool {
	if s == nil {
		return false
	}
	s.RLock()
	defer s.RUnlock()
	if s.master != nil {
		return s.master.UUID == UUID
	}
	return false
}

func (s *server) resetServer() {
	s.Lock()
	defer s.Unlock()
	s.servers = make([]*server_info.ServerInfo, 0)
	s.master = nil
}

// setServersChan 当监听到服务节点变化时，将最新的服务节点信息放入该channel里
func (s *server) setServersChan() {
	// 即使没有其他服务消费该channel，也能保证该channel不会阻塞
	for len(s.serversChan) >= 1 {
		<-s.serversChan
	}
	s.serversChan <- s.getInstances()
}

func (s *server) run() {
	log.Infof("start to discover cc component from zk, path:[%s].", s.path)
	go func() {
		for svr := range s.discoverChan {
			log.Warnf("received one zk event from path %s.", s.path)
			if svr.Err != nil {
				log.Errorf("get zk event with error about path[%s]. err: %v", s.path, svr.Err)
				continue
			}

			if len(svr.Server) <= 0 {
				log.Warnf("get zk event with 0 instance with path[%s], reset its servers", s.path)
				s.resetServer()
				s.setServersChan()
				continue
			}

			s.updateServer(svr.Server)
			s.setServersChan()
		}
	}()
}

func (s *server) updateServer(svrs []string) {
	servers := make([]*server_info.ServerInfo, 0)
	var master *server_info.ServerInfo

	for _, svr := range svrs {
		server := new(server_info.ServerInfo)
		if err := json.Unmarshal([]byte(svr), server); err != nil {
			log.Errorf("unmarshal server info failed, zk path[%s], err: %v", s.path, err)
			continue
		}

		if server.Scheme != "https" {
			server.Scheme = "http"
		}

		if server.Port == 0 {
			log.Errorf("invalid port 0, with zk path: %s", s.path)
			continue
		}

		if len(server.RegisterIP) == 0 {
			log.Errorf("invalid ip with zk path: %s", s.path)
			continue
		}

		if server.Environment == s.environment {
			servers = append(servers, server)
		}
		if master == nil {
			master = server
		}
	}

	s.Lock()
	s.servers = servers
	s.master = master
	s.Unlock()

	log.Infof("update component with new server instance %s about path: %s", servers, s.path)

}
