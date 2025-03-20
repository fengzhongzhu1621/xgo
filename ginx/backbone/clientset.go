package backbone

import (
	"github.com/fengzhongzhu1621/xgo/collections/flowctrl"
	"github.com/fengzhongzhu1621/xgo/ginx/discovery"
	"github.com/fengzhongzhu1621/xgo/ginx/rest"
	"github.com/fengzhongzhu1621/xgo/ginx/service/coreservice"
	"github.com/fengzhongzhu1621/xgo/ginx/types"
)

type ClientSet struct {
	version  string
	client   rest.IHttpClient
	discover discovery.IDiscoveryInterface
	throttle flowctrl.RateLimiter
	Mock     types.MockInfo
}

func (cs *ClientSet) CoreService() coreservice.CoreServiceClientInterface {
	c := &rest.Capability{
		Client:   cs.client,
		Discover: cs.discover.CoreService(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	return coreservice.NewCoreServiceClient(c, cs.version)
}

func NewMockClientSet() *ClientSet {
	return &ClientSet{
		version:  "unit_test",
		client:   nil,
		discover: discovery.NewMockDiscoveryInterface(),
		throttle: flowctrl.NewMockRateLimiter(),
		Mock:     types.MockInfo{Mocked: true},
	}
}

func NewClientSet(client rest.IHttpClient, discover discovery.IDiscoveryInterface,
	throttle flowctrl.RateLimiter) IClientSetInterface {
	return &ClientSet{
		version:  "v3",
		client:   client,
		discover: discover,
		throttle: throttle,
	}
}
