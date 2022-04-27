package request

import (
	opts "xgo/services/name/client/options"
)

// 客户端请求接口
type IRequest interface {
	// Init 初始化
	Init(o opts.InitOptions) (err error)
}

// var (
//	_ IRequest     = HttpRequest{}
//)
