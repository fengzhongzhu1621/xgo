package webserver

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/backbone/configcenter"
)

type WebServer struct {
	Config config.WebServerConfig
}

func (w *WebServer) onServerConfigUpdate(previous, current configcenter.ProcessConfig) {

}

func (w *WebServer) Stop() error {
	return nil
}
