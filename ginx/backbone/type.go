package backbone

import (
	"net/http"

	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
	"github.com/fengzhongzhu1621/xgo/ginx/service/coreservice"
	"github.com/fengzhongzhu1621/xgo/ginx/utils/tls"
	"github.com/fengzhongzhu1621/xgo/network/ssl"
)

type IClientSetInterface interface {
	CoreService() coreservice.CoreServiceClientInterface
}

type Config struct {
	RegisterPath string
	RegisterInfo server_info.ServerInfo
	CoreAPI      IClientSetInterface
}

// Server TODO
type Server struct {
	ListenAddr   string
	ListenPort   uint
	Handler      http.Handler
	TLS          *ssl.TLSClientConfig
	PProfEnabled bool
}

type APIMachineryConfig struct {
	// request's qps value
	QPS int64
	// request's burst value
	Burst     int64
	TLSConfig *ssl.TLSClientConfig
	ExtraConf *tls.ExtraClientConfig
}
