package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/gin-gonic/gin"
)

// 热重载配置: 单独监听 SIGHUP，不重启进程
// 优雅关闭: signal.NotifyContext + server.Shutdown()
// systemd Watchdog 自动重启: daemon.SdNotify(false, "WATCHDOG=1")
//

func main() {
	// 创建一个 context.Context，当收到 SIGINT 或 SIGTERM 时，ctx.Done() 会返回一个关闭的 channel。
	// 返回的 stop 函数用于取消信号监听（避免内存泄漏）。
	// SIGINT（Ctrl+C）和 SIGTERM（systemd 默认终止信号）
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// 确保在程序退出前释放信号监听资源。
	defer stop()

	// 初始化 Gin 路由
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(10 * time.Second) // 模拟长请求
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	// 配置 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 启动 systemd Watchdog 心跳（需在 systemd 单元中设置 WatchdogSec=10）
	go func() {
		ticker := time.NewTicker(5 * time.Second) // 间隔必须小于 WatchdogSec
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if _, err := daemon.SdNotify(false, "WATCHDOG=1"); err != nil {
					log.Printf("Failed to send watchdog heartbeat: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// 单独监听 SIGHUP 信号（热重载配置）
	go func() {
		sighupChan := make(chan os.Signal, 1)
		signal.Notify(sighupChan, syscall.SIGHUP)
		for {
			select {
			case <-sighupChan:
				log.Println("SIGHUP received, reloading config...")
				// TODO 在这里重新加载配置文件（例如重新读取环境变量或配置文件）
			case <-ctx.Done():
				return
			}
		}
	}()

	// 在 goroutine 中启动服务器（避免阻塞优雅关闭逻辑）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待终止信号
	<-ctx.Done()

	// 恢复默认信号行为并打印关闭日志
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// 设置 5 秒超时强制关闭
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
