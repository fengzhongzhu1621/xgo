package backbone

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fengzhongzhu1621/xgo/collections/flowctrl"
	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/db/zookeeper/confregdiscover"
	"github.com/fengzhongzhu1621/xgo/db/zookeeper/registerdiscover"
	"github.com/fengzhongzhu1621/xgo/ginx/backbone/configcenter"
	"github.com/fengzhongzhu1621/xgo/ginx/discovery"
	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
	"github.com/fengzhongzhu1621/xgo/monitor"
	"github.com/fengzhongzhu1621/xgo/monitor/metrics"
	"github.com/fengzhongzhu1621/xgo/monitor/opentelemetry"
	log "github.com/sirupsen/logrus"
)

// BackboneParameter Used to constrain different services to ensure
// consistency of service startup capabilities
type BackboneParameter struct {
	// ConfigUpdate handle process config change
	ConfigUpdate configcenter.ProcHandlerFunc
	ExtraUpdate  configcenter.ProcHandlerFunc
	// config path
	ConfigPath string
	// http server parameter
	SrvInfo *server_info.ServerInfo
	SrvRegdiscv
}

func NewBackbone(ctx context.Context, input *BackboneParameter) (*Engine, error) {
	// if err := validateParameter(input); err != nil {
	// 	return nil, err
	// }

	metricService := metrics.NewService(metrics.Config{ProcessName: server_option.GetIdentification(),
		ProcessInstance: input.SrvInfo.Instance()})

	server_info.SetServerInfo(input.SrvInfo)

	engine, err := New()
	if err != nil {
		return nil, fmt.Errorf("new engine failed, err: %v", err)
	}
	engine.registerPath = getRegisterPath(input.SrvInfo.IP)
	engine.srvInfo = input.SrvInfo
	engine.metric = metricService
	engine.Disable = input.Disable

	handler := &configcenter.CCHandler{
		// 扩展这个函数， 新加传递错误
		OnProcessUpdate: input.ConfigUpdate,
		OnExtraUpdate:   input.ExtraUpdate,
		OnMongodbUpdate: engine.onMongodbUpdate,
		OnRedisUpdate:   engine.onRedisUpdate,
	}

	if !input.Disable {
		client, err := registerdiscover.NewSvcManagerClient(ctx, input.Regdiscv)
		if err != nil {
			return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
		}
		serviceDiscovery, err := discovery.NewServiceDiscovery(client, input.SrvInfo.Environment)
		if err != nil {
			return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
		}
		disc, err := NewServiceRegister(client)
		if err != nil {
			return nil, fmt.Errorf("new service discover failed, err:%v", err)
		}

		engine.client = client
		engine.discovery = serviceDiscovery
		engine.ServiceManageInterface = serviceDiscovery
		engine.SvcDisc = disc

		// add default configcenter
		zkdisc := confregdiscover.NewZkRegDiscover(client)
		configCenter := &configcenter.ConfigCenter{
			Type:               configcenter.DefaultConfigCenter,
			ConfigCenterDetail: zkdisc,
		}
		configcenter.AddConfigCenter(configCenter)

		tlsConf, err := GetTLSConf()
		if err != nil {
			log.Errorf("get tls config error, err: %v", err)
			return nil, err
		}
		engine.apiMachineryConfig = &APIMachineryConfig{
			QPS:       1000,
			Burst:     2000,
			TLSConfig: tlsConf,
		}

		machinery, err := newApiMachinery(serviceDiscovery, engine.apiMachineryConfig)
		if err != nil {
			return nil, err
		}
		engine.CoreAPI = machinery

		// if err = handleNotice(ctx, client.Client(), input.SrvInfo.Instance()); err != nil {
		// 	return nil, fmt.Errorf("handle notice failed, err: %v", err)
		// }
	}

	// get the real configuration center.
	currentConfigCenter := configcenter.CurrentConfigCenter()

	if err = configcenter.NewConfigCenter(ctx, currentConfigCenter, input.ConfigPath, handler); err != nil {
		return nil, fmt.Errorf("new config center failed, err: %v", err)
	}

	if err := monitor.InitMonitor(); err != nil {
		return nil, fmt.Errorf("init monitor failed, err: %v", err)
	}

	if err := opentelemetry.InitOpenTelemetryConfig(); err != nil {
		return nil, fmt.Errorf("init openTelemetry config failed, err: %v", err)
	}

	if err := opentelemetry.InitTracer(ctx); err != nil {
		return nil, fmt.Errorf("init tracer failed, err: %v", err)
	}

	return engine, nil
}

func NewApiMachinery(c *APIMachineryConfig, discover discovery.IDiscoveryInterface) (IClientSetInterface, error) {
	extraConf := make([]ExtraClientConfig, 0)
	if c.ExtraConf != nil {
		extraConf = append(extraConf, *c.ExtraConf)
	}
	client, err := NewClient(c.TLSConfig, extraConf...)
	if err != nil {
		return nil, err
	}

	flowcontrol := flowctrl.NewRateLimiter(c.QPS, c.Burst)
	return NewClientSet(client, discover, flowcontrol), nil
}

func newApiMachinery(disc discovery.IDiscoveryInterface,
	config *APIMachineryConfig) (IClientSetInterface, error) {

	machinery, err := NewApiMachinery(config, disc)
	if err != nil {
		return nil, fmt.Errorf("new api machinery failed, err: %v", err)
	}

	return machinery, nil
}

func getRegisterPath(ip string) string {
	return fmt.Sprintf("%s/%s/%s", configcenter.SERV_BASEPATH, server_option.GetIdentification(), ip)
}

func StartServer(ctx context.Context, cancel context.CancelFunc, e *Engine, HTTPHandler http.Handler,
	pprofEnabled bool) error {
	tlsConf, err := GetTLSConf()
	if err != nil {
		log.Errorf("get tls config error, err: %v", err)
		return err
	}

	if IsTLS(tlsConf) {
		e.srvInfo.Scheme = "https"
	}

	e.server = Server{
		ListenAddr:   e.srvInfo.IP,
		ListenPort:   e.srvInfo.Port,
		Handler:      e.Metric().HTTPMiddleware(HTTPHandler),
		TLS:          tlsConf,
		PProfEnabled: pprofEnabled,
	}

	if err := ListenAndServe(e.server, e.SvcDisc, cancel); err != nil {
		return err
	}

	// wait for a while to see if ListenAndServe in goroutine is successful
	// to avoid registering an invalid server address on zk
	time.Sleep(time.Second)

	if e.Disable {
		return nil
	}

	return e.SvcDisc.Register(e.registerPath, *e.srvInfo)
}
