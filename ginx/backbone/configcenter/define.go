package configcenter

const (
	DefaultConfigCenter = "zookeeper"
	// connect svcManager retry connect time
	maxRetry = 200
)

const (
	SERV_BASEPATH     = "/services/endpoints"
	SERVCONF_BASEPATH = "/services/config"
)

const (
	ConfigureRedis  = "redis"
	ConfigureMongo  = "mongodb"
	ConfigureCommon = "common"
	ConfigureExtra  = "extra"
)
