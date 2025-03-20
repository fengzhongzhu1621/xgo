package webserver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/ginx/backbone"
	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
	"github.com/fengzhongzhu1621/xgo/ginx/service/webservice"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *server_option.ServerOption) error {
	// 根据命令行参数创建服务信息对象
	svrInfo, err := server_info.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	webSvr := new(WebServer)
	webSvr.Config.DeploymentMethod = op.DeploymentMethod

	input := &backbone.BackboneParameter{
		ConfigUpdate: webSvr.onServerConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		SrvRegdiscv:  backbone.SrvRegdiscv{Regdiscv: op.ServConf.RegDiscover},
		SrvInfo:      svrInfo,
	}

	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	configReady := false
	for sleepCnt := 0; sleepCnt < constant.APPConfigWaitTime; sleepCnt++ {
		if "" == webSvr.Config.Site.DomainUrl {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if !configReady {
		return errors.New("configuration item not found")
	}

	service, err := initWebService(webSvr, engine)
	if err != nil {
		return err
	}

	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), false)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	}

	return nil
}

func initWebService(webSvr *WebServer, engine *backbone.Engine) (*webservice.Service, error) {
	service := new(webservice.Service)

	return service, nil
}
