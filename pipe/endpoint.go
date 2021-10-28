package pipe

type Endpoint interface {
	StartReading() // 启动协程读取数据，写入到管道中
	Terminate()
	Output() chan []byte // 返回管道
	Send([]byte) bool
}

// 双工管道传递交换数据
func PipeEndpoints(e1, e2 Endpoint) {
	// 启动协程读取数据，写入到管道中
	e1.StartReading()
	e2.StartReading()

	defer e1.Terminate()
	defer e2.Terminate()

	// 双向管道数据发送
	for {
		select {
		case msgOne, ok := <-e1.Output():
			if !ok || !e2.Send(msgOne) {
				return
			}
		case msgTwo, ok := <-e2.Output():
			if !ok || !e1.Send(msgTwo) {
				return
			}
		}
	}
}
