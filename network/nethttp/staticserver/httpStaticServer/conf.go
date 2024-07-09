package httpStaticServer

import (
	"github.com/fengzhongzhu1621/xgo/network/nethttp/staticServer"
)

// 构造默认配置对象
func (s *HTTPStaticServer) defaultAccessConf() staticServer.AccessConf {
	return staticServer.AccessConf{
		Upload: s.Upload,
		Delete: s.Delete,
	}
}
