package backbone

import (
	"github.com/fengzhongzhu1621/xgo/db/zookeeper/registerdiscover"
	"github.com/fengzhongzhu1621/xgo/ginx/discovery"
)

type SrvRegdiscv struct {
	client                 *registerdiscover.ZkClient
	ServiceManageInterface discovery.IServiceManageInterface
	SvcDisc                ServiceRegisterInterface
	discovery              discovery.IDiscoveryInterface
	// registerPath the path registered to the Service Discovery Center
	registerPath string
	// service component addr
	Regdiscv string
	// Disable disable service registration discovery
	Disable bool
}

// Discovery return discovery
func (s *SrvRegdiscv) Discovery() discovery.IDiscoveryInterface {
	return s.discovery
}

// ServiceManageClient return service manage client
func (s *SrvRegdiscv) ServiceManageClient() *registerdiscover.ZkClient {
	return s.client
}
