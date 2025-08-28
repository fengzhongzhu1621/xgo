package main

import (
	"log"
	"net/http"
	// 启用 pprof 性能分析
	// 向 net/http 的默认 server——DefaultServerMux 中完成了一系列路径及对应 handler 的注册
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/wolfogre/go-pprof-practice/animal"
)

func main() {
	// 全称 "Go Maximum Procs"（Go 最大处理器数），控制 Go 程序可以并行执行的线程数量。
	// 参数 1：表示将最大并行线程数限制为 （即单线程模式）。
	runtime.GOMAXPROCS(1)
	// 启用 mutex 性能分析
	runtime.SetMutexProfileFraction(1)
	// 启用 block 性能分析
	runtime.SetBlockProfileRate(1)

	go func() {
		// 启动 http server. 对应 pprof 的一系列 handler 也会挂载在该端口下
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	// 运行各项动物的活动
	for {
		for _, v := range animal.AllAnimals {
			v.Live()
		}
		time.Sleep(time.Second)
	}
}
