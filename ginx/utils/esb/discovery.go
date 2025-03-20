package esb

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// EsbConfigSrv TODO
type EsbConfigSrv struct {
	addrs     string
	appCode   string
	appSecret string
	sync.RWMutex
}

// EsbSrvDiscoveryInterface TODO
type EsbSrvDiscoveryInterface interface {
	GetServers() ([]string, error)
}

// NewEsbConfigSrv TODO
func NewEsbConfigSrv(srvChan chan EsbConfig, defaultCfg *EsbConfig) *EsbConfigSrv {
	esb := &EsbConfigSrv{}

	if defaultCfg != nil {
		esb.addrs = defaultCfg.Addrs
		esb.appCode = defaultCfg.AppCode
		esb.appSecret = defaultCfg.AppSecret
	}

	go func() {
		if srvChan == nil {
			return
		}
		for {
			config := <-srvChan
			esb.Lock()
			esb.addrs = config.Addrs
			esb.appCode = config.AppCode
			esb.appSecret = config.AppSecret
			log.Infof("cmdb config changed, config: %+v", config)
			esb.Unlock()
		}
	}()

	return esb
}

// GetEsbSrvDiscoveryInterface TODO
func (esb *EsbConfigSrv) GetEsbSrvDiscoveryInterface() EsbSrvDiscoveryInterface {
	// maybe will deal some logic about server
	return esb
}

// GetServers TODO
func (esb *EsbConfigSrv) GetServers() ([]string, error) {
	// maybe will deal some logic about server
	esb.RLock()
	defer esb.RUnlock()
	return []string{esb.addrs}, nil
}

// GetServersChan TODO
func (esb *EsbConfigSrv) GetServersChan() chan []string {
	return nil
}

// GetConfig TODO
func (esb *EsbConfigSrv) GetConfig() EsbConfig {
	esb.RLock()
	defer esb.RUnlock()
	return EsbConfig{
		Addrs:     esb.addrs,
		AppCode:   esb.appCode,
		AppSecret: esb.appSecret,
	}
}
