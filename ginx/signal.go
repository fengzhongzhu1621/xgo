package ginx

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Interrupt 监听操作系统发送的信号，并在接收到特定信号时执行传入的回调函数 onSignal
func Interrupt(onSignal func()) {
	// 创建一个信号通道，用于接收操作系统的信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// 等待并处理信号
	for s := range c {
		log.Infof("Caught signal %s. Exiting.", s)
		// 调用传入的 onSignal 回调函数
		onSignal()
		// 关闭通道 c，以确保循环能够退出
		close(c)
	}
}
