package registerdiscover

// RegDiscvServer define the register and discover function interface
type RegDiscvServer interface {
	// Ping to ping server
	Ping() error
	// RegisterAndWatch register server info into register-discover service platform,
	// and watch the info, if not exist, then register again
	RegisterAndWatch(key string, data []byte) error
	// GetServNodes get server nodes
	GetServNodes(key string) ([]string, error)
	// Discover server from the register-discover service platform
	Discover(key string) (<-chan *DiscoverEvent, error)
	// Cancel to stop server register and discover
	Cancel()
	// ClearRegisterPath to delete server register path from zk
	ClearRegisterPath() error
}
