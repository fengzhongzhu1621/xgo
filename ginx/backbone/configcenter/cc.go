package configcenter

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo/config/viper/viper_parser"
	"github.com/fengzhongzhu1621/xgo/db/zookeeper/confregdiscover"
	log "github.com/sirupsen/logrus"
)

var confC *CC

type ProcHandlerFunc func(previous, current ProcessConfig)

type CCHandler struct {
	OnProcessUpdate ProcHandlerFunc
	OnExtraUpdate   ProcHandlerFunc
	OnMongodbUpdate func(previous, current ProcessConfig)
	OnRedisUpdate   func(previous, current ProcessConfig)
}

type CC struct {
	sync.Mutex
	// used to stop the config center gracefully.
	ctx             context.Context
	disc            confregdiscover.ConfRegDiscvIf
	handler         *CCHandler
	procName        string
	previousProc    *ProcessConfig
	previousExtra   *ProcessConfig
	previousMongodb *ProcessConfig
	previousRedis   *ProcessConfig
}

func NewConfigCenter(ctx context.Context, disc confregdiscover.ConfRegDiscvIf, confPath string, handler *CCHandler) error {
	return New(ctx, confPath, disc, handler)
}

func New(ctx context.Context, confPath string, disc confregdiscover.ConfRegDiscvIf, handler *CCHandler) error {
	confC = &CC{
		ctx:          ctx,
		disc:         disc,
		handler:      handler,
		previousProc: new(ProcessConfig),
	}

	// parse config only from file
	if len(confPath) != 0 {
		return GetLocalConf(confPath, handler)
	}

	if err := confC.run(); err != nil {
		return err
	}

	confC.sync()

	return nil
}

func (c *CC) run() error {
	commonConfPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureCommon)
	commonConfEvent, err := c.disc.Discover(commonConfPath)
	if err != nil {
		return err
	}

	extraConfPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureExtra)
	extraConfEvent, err := c.disc.Discover(extraConfPath)
	if err != nil {
		return err
	}

	mongodbConfPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureMongo)
	mongodbConfEvent, err := c.disc.Discover(mongodbConfPath)
	if err != nil {
		return err
	}

	redisConfPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureRedis)
	redisConfEvent, err := c.disc.Discover(redisConfPath)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case pEvent := <-commonConfEvent:
				c.onProcChange(pEvent)
			case pEvent := <-extraConfEvent:
				c.onExtraChange(pEvent)
			case pEvent := <-mongodbConfEvent:
				c.onMongodbChange(pEvent)
			case pEvent := <-redisConfEvent:
				c.onRedisChange(pEvent)
			case <-c.ctx.Done():
				log.Warnf("config center event watch stopped because of context done.")
				return
			}
		}
	}()
	return nil
}

func (c *CC) onProcChange(cur *confregdiscover.DiscoverEvent) {
	if cur.Err != nil {
		log.Errorf("config center received event that %s config has changed, but got err: %v", ConfigureCommon,
			cur.Err)
		return
	}

	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousProc
	c.previousProc = now
	if err := viper_parser.SetCommonFromByte(now.ConfigData); err != nil {
		log.Errorf("add updated configuration error: %v", err)
		return
	}
	if c.handler != nil {
		if c.handler.OnProcessUpdate != nil {
			go c.handler.OnProcessUpdate(*prev, *now)
		}
	}
}

func (c *CC) onExtraChange(cur *confregdiscover.DiscoverEvent) {
	if cur.Err != nil {
		log.Errorf("config center received event that %s config has changed, but got err: %v", ConfigureExtra,
			cur.Err)
		return
	}

	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousExtra
	if prev == nil {
		prev = &ProcessConfig{}
	}
	c.previousExtra = now
	if err := viper_parser.SetExtraFromByte(now.ConfigData); err != nil {
		log.Errorf("add updated extra configuration error: %v", err)
		return
	}
	if c.handler != nil {
		if c.handler.OnExtraUpdate != nil {
			go c.handler.OnExtraUpdate(*prev, *now)
		}
	}
}

func (c *CC) onMongodbChange(cur *confregdiscover.DiscoverEvent) {
	if cur.Err != nil {
		log.Errorf("config center received event that %s config has changed, but got err: %v", ConfigureCommon,
			cur.Err)
		return
	}
	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousMongodb
	if prev == nil {
		prev = &ProcessConfig{}
	}
	c.previousMongodb = now
	if c.handler != nil {
		if c.handler.OnMongodbUpdate != nil {
			go c.handler.OnMongodbUpdate(*prev, *now)
		}
	}
}

func (c *CC) onRedisChange(cur *confregdiscover.DiscoverEvent) {
	if cur.Err != nil {
		log.Errorf("config center received event that %s config has changed, but got err: %v", ConfigureCommon,
			cur.Err)
		return
	}
	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousRedis
	if prev == nil {
		prev = &ProcessConfig{}
	}
	c.previousRedis = now
	if c.handler != nil {
		if c.handler.OnRedisUpdate != nil {
			go c.handler.OnRedisUpdate(*prev, *now)
		}
	}
}

func (c *CC) sync() {
	log.Infof("start sync config from config center.")
	c.syncProc()
	c.syncExtra()
	c.syncMongodb()
	c.syncRedis()
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:

			}
			// sync the data from zk, and compare if it has been changed.
			// then call their handler.
			c.syncProc()
			c.syncExtra()
			c.syncMongodb()
			c.syncRedis()
			time.Sleep(15 * time.Second)
		}
	}()
}

func (c *CC) syncProc() {
	log.Infof("start sync proc config from config center.")
	procPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureCommon)
	data, err := c.disc.Read(procPath)
	if err != nil {
		log.Errorf("sync process config failed, node: %s, err: %v", procPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))

	c.Lock()
	if reflect.DeepEqual(conf, c.previousProc) {
		log.Infof("sync process config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &confregdiscover.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()
	c.onProcChange(event)
}

func (c *CC) syncExtra() {
	log.Infof("start sync proc config from config center.")
	procPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureExtra)
	data, err := c.disc.Read(procPath)
	if err != nil {
		log.Errorf("sync *extra* config failed, node: %s, err: %v", procPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))

	c.Lock()
	if reflect.DeepEqual(conf, c.previousExtra) {
		log.Infof("sync *extra* config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &confregdiscover.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()

	c.onExtraChange(event)

}

func (c *CC) syncMongodb() {
	log.Infof("start sync mongo config from config center.")
	mongoPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureMongo)
	data, err := c.disc.Read(mongoPath)
	if err != nil {
		log.Errorf("sync *mongo* config failed, node: %s, err: %v", mongoPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))
	c.Lock()
	if reflect.DeepEqual(conf, c.previousMongodb) {
		log.Infof("sync *mongo* config, but nothing is changed.")
		c.Unlock()
		return
	}
	event := &confregdiscover.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()

	c.onMongodbChange(event)
}

// GetLocalConf get local config
func GetLocalConf(confPath string, handler *CCHandler) error {
	if err := viper_parser.SetLocalFile(confPath); err != nil {
		return fmt.Errorf("parse config file: %s failed, err: %v", confPath, err)
	}

	if handler.OnProcessUpdate != nil {
		handler.OnProcessUpdate(ProcessConfig{}, ProcessConfig{})
	}

	return nil
}

func (c *CC) syncRedis() {
	log.Infof("start sync redis config from config center.")
	redisPath := fmt.Sprintf("%s/%s", SERVCONF_BASEPATH, ConfigureRedis)
	data, err := c.disc.Read(redisPath)
	if err != nil {
		log.Errorf("sync *redis* config failed, node: %s, err: %v", redisPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))

	c.Lock()
	if reflect.DeepEqual(conf, c.previousRedis) {
		log.Infof("sync *redis* config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &confregdiscover.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()

	c.onRedisChange(event)

}
