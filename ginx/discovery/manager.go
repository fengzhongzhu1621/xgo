package discovery

var (
	// needDiscoveryServiceName 服务依赖的第三方服务名字的配置
	needDiscoveryServiceName = make(map[string]struct{}, 0)
)

// DiscoveryAllService 发现所有定义的服务
func DiscoveryAllService() {
	for name := range AllModule {
		needDiscoveryServiceName[name] = struct{}{}
	}
}

// AddDiscoveryService 新加需要发现服务的名字
func AddDiscoveryService(name ...string) {
	for _, name := range name {
		needDiscoveryServiceName[name] = struct{}{}
	}
}

// GetDiscoveryService TODO
func GetDiscoveryService() map[string]struct{} {
	// compatible 如果没配置,发现所有的服务
	if len(needDiscoveryServiceName) == 0 {
		DiscoveryAllService()
	}

	return needDiscoveryServiceName
}
