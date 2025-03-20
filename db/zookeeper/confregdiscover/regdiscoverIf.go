package confregdiscover

// ConfRegDiscvIf define the configurate register and discover function interface
type ConfRegDiscvIf interface {
	// Ping to ping server
	Ping() error
	// Write the config data into configure register-discover service
	Write(key string, data []byte) error
	// Read the config data from configure register-discover service
	Read(key string) (string, error)
	// Discover the config change
	Discover(key string) (<-chan *DiscoverEvent, error)
}
