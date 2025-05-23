package registerdiscover

import (
	"time"
)

// DiscoverEvent if servers chenged, will create a discover event
type DiscoverEvent struct { //
	Err    error
	Key    string
	Server []string
	Nodes  []string
}

// RegDiscover is data struct of register-discover
type RegDiscover struct {
	rdServer RegDiscvServer
}

// NewRegDiscover used to create a object of RegDiscover
func NewRegDiscover(client *ZkClient, timeout time.Duration) *RegDiscover {
	regDiscv := &RegDiscover{
		rdServer: nil,
	}

	regDiscv.rdServer = RegDiscvServer(NewZkRegDiscv(client))

	return regDiscv
}

// NewRegDiscoverEx used to create a object of RegDiscover
func NewRegDiscoverEx(client *ZkClient) *RegDiscover {
	regDiscv := &RegDiscover{
		rdServer: nil,
	}

	regDiscv.rdServer = RegDiscvServer(NewZkRegDiscv(client))

	return regDiscv
}

// RegisterAndWatchService register service info into register-discover platform
// and then watch the service info, if not exist, then register again
// key is the index of registered service
// data is the service information
func (rd *RegDiscover) RegisterAndWatchService(key string, data []byte) error {
	return rd.rdServer.RegisterAndWatch(key, data)
}

// GetServNodes TODO
func (rd *RegDiscover) GetServNodes(key string) ([]string, error) {
	return rd.rdServer.GetServNodes(key)
}

// DiscoverService used to discover the service that registered in `key`
func (rd *RegDiscover) DiscoverService(key string) (<-chan *DiscoverEvent, error) {
	return rd.rdServer.Discover(key)
}

// Ping to ping server
func (rd *RegDiscover) Ping() error {
	return rd.rdServer.Ping()
}

// Cancel to stop server register and discover
func (rd *RegDiscover) Cancel() {
	rd.rdServer.Cancel()
}

// ClearRegisterPath to delete server register path from zk
func (rd *RegDiscover) ClearRegisterPath() error {
	return rd.rdServer.ClearRegisterPath()
}
