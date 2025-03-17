package backbone

import (
	"net/http"

	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
)

type ClientSetInterface interface {
}

type Config struct {
	RegisterPath string
	RegisterInfo server_info.ServerInfo
	CoreAPI      ClientSetInterface
}

// Server TODO
type Server struct {
	ListenAddr   string
	ListenPort   uint
	Handler      http.Handler
	TLS          *TLSClientConfig
	PProfEnabled bool
}
