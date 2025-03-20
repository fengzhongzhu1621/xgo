package backbone

import (
	"sync"

	"github.com/fengzhongzhu1621/xgo/channel"
	"github.com/fengzhongzhu1621/xgo/config/viper/viper_parser"
	"github.com/fengzhongzhu1621/xgo/db/mongodb"
	redisClient "github.com/fengzhongzhu1621/xgo/db/redis/client"
	"github.com/fengzhongzhu1621/xgo/ginx/backbone/configcenter"
	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
	"github.com/fengzhongzhu1621/xgo/monitor/metrics"
	log "github.com/sirupsen/logrus"
)

// New new engine
func New() (*Engine, error) {
	return &Engine{
		ctx: channel.NewContext(),
	}, nil
}

// Engine TODO
type Engine struct {
	CoreAPI            IClientSetInterface
	apiMachineryConfig *APIMachineryConfig
	metric             *metrics.Service
	sync.Mutex
	server  Server
	srvInfo *server_info.ServerInfo
	ctx     channel.IContextInterface
	SrvRegdiscv
}

// ApiMachineryConfig TODO
func (e *Engine) ApiMachineryConfig() *APIMachineryConfig {
	return e.apiMachineryConfig
}

// Metric TODO
func (e *Engine) Metric() *metrics.Service {
	return e.metric
}

// GetSrvInfo get service info
func (e *Engine) GetSrvInfo() *server_info.ServerInfo {
	return e.srvInfo
}

func (e *Engine) Ping() error {
	if e.SrvRegdiscv.Disable {
		return nil
	}
	return e.SvcDisc.Ping()
}

func (e *Engine) onMongodbUpdate(previous, current configcenter.ProcessConfig) {
	e.Lock()
	defer e.Unlock()
	if err := viper_parser.SetMongodbFromByte(current.ConfigData); err != nil {
		log.Errorf("parse mongo config failed, err: %s, data: %s", err.Error(), string(current.ConfigData))
	}
}

func (e *Engine) onRedisUpdate(previous, current configcenter.ProcessConfig) {
	e.Lock()
	defer e.Unlock()
	if err := viper_parser.SetRedisFromByte(current.ConfigData); err != nil {
		log.Errorf("parse redis config failed, err: %s, data: %s", err.Error(), string(current.ConfigData))
	}
}

func (e *Engine) WithRedis(waitTime int, prefixes ...string) (redisClient.Config, error) {
	// use default prefix if no prefix is specified, or use the first prefix
	var prefix string
	if len(prefixes) == 0 {
		prefix = "redis"
	} else {
		prefix = prefixes[0]
	}

	return viper_parser.Redis(prefix, waitTime)
}

func (e *Engine) WithMongo(waitTime int, prefixes ...string) (mongodb.Config, error) {
	var prefix string
	if len(prefixes) == 0 {
		prefix = "mongodb"
	} else {
		prefix = prefixes[0]
	}

	return viper_parser.Mongo(prefix, waitTime)
}
