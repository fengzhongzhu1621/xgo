package ginx

import (
	"context"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/server"
)

func Init() {

}

func Start() {
	Init()

	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		// 监听操作系统发送的信号，并在接收到特定信号时取消上下文
		Interrupt(cancelFunc)
	}()

	globalConfig := config.GetGlobalConfig()
	httpServer := server.NewServer(globalConfig)

	httpServer.Run(ctx)
}
