package webserver

import (
	"context"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/config/server_option"
)

type WebServer struct {
	Config config.WebServerConfig
}

func Run(ctx context.Context, cancel context.CancelFunc, op *server_option.ServerOption) error {

	webSvr := new(WebServer)
	webSvr.Config.DeploymentMethod = op.DeploymentMethod

	select {
	case <-ctx.Done():
	}

	return nil
}
