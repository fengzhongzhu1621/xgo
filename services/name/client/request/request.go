package request

import (
	"github.com/fengzhongzhu1621/xgo/services/name/client/options"
)

/**
 * 客户端请求接口
 */
type IRequest interface {
	// Init 初始化
	Init(o options.InitOptions) (err error)
}

//var (
//	_ IRequest     = HttpRequest{}
//)
