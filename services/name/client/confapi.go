package name

import (
	"fmt"
	opts "github.com/fengzhongzhu1621/xgo/services/name/client/options"
	request "github.com/fengzhongzhu1621/xgo/services/name/client/request"
)


type ConfApi struct {
	options     opts.InitOptions
	handler		request.IRequest
}

func (api *ConfApi) Init() error {
	// 根据协议类型，创建不同的连接器
	if api.options.ProtoType == "http" {
		api.handler = &request.HttpRequest{}
	} else {
		return fmt.Errorf("proto type %s not supported", api.options.ProtoType)
	}
	return api.handler.Init(api.options)
}

