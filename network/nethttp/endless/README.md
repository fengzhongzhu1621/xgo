# 简介
优雅的零停机重启库endless

* 平滑套接字交接：通过SO_REUSEPORT实现端口复用
* 双进程协作：新旧进程并行运行直至旧连接完成
* 信号驱动：支持SIGHUP等信号触发安全重启


```go
              [旧进程]
                │
接收SIGHUP信号───┤
                ├─► 创建新进程（继承文件描述符）
                │
          [新旧进程共存期]──┬─► 新连接路由到新进程
                          └─► 旧进程处理存量请求
```

# ListenAndServe

```go
package main
import (
    "net/http"
    "github.com/fvbock/endless"
)
func handler(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Hello Zero-Downtime!"))
}

func main() {
    server := endless.NewServer(":8080", http.HandlerFunc(handler))
    server.BeforeBegin = func(addr string) {
        log.Printf("Listening on %s", addr)
    }
    server.ListenAndServe()
}
```

```go
func main() {
    router := gin.Default()
    router.GET("/", func(c *gin.Context)
        c.String(200,"Gin with endless!")
    })

    server := endless.NewServer(":8080", router)
    server.ReadTimeout = 15* time.Second
    server.WriteTimeout = 30* time.Second
    if err := server.ListenAndServe(); err !=nil{
        log.Fatal("Server error:", err)
    }}
```

# 自定义信号处理
```go
server := endless.NewServer(...)

// 添加自定义信号处理器
server.SignalHooks[endless.PRE_SIGNAL] = append(
    server.SignalHooks[endless.PRE_SIGNAL],
    func() {
        log.Println("Preparing for restart...")
        // 执行预清理操作
    },
)
```

# 优雅关闭
```go
server := endless.NewServer(...)

// 设置优雅关闭超时
server.ShutdownInitiated = func() {
    log.Println("Starting graceful shutdown")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 关闭后台任务
    closeBackgroundWorkers()

    <-ctx.Done()
}
```

# systemd 服务
```
[Unit]
Description=My Custom Service
Restart=always
RestartSec=3

[Service]
Type=notify
ExecStart=/usr/bin/my-service
ExecReload=/bin/kill -HUP $MAINPID
WatchdogSec=10
NotifyAccess=all

[Install]
WantedBy=multi-user.target
```
Unit
* Restart
定义了当服务以任何原因退出时，systemd 是否以及如何自动重启服务。always 表示无论服务是如何退出的（无论是正常退出还是由于错误退出），systemd 都会尝试重新启动服务。其他常见的选项包括 on-failure（仅在非正常退出时重启）、no（从不重启）等
* RestartSec
  指定了在服务退出后，systemd 等待多长时间再尝试重启服务。单位默认为秒，所以这里表示等待 3 秒后重启。这可以防止服务在短时间内频繁重启，有助于避免潜在的资源耗尽问题。

Service
* Type=notify 用于启用 Watchdog 功能。这要求你的服务能够发送通知给 systemd。
* ExecStart 指定了启动服务的命令。确保 /usr/bin/my-service 存在并且具有可执行权限。
* ExecReload 定义了当服务接收到 reload 信号时，systemd 将执行的命令。确保你的服务确实支持并正确处理 SIGHUP 信号，否则这种重新加载方式可能无效或导致不可预期的行为。
  - -HUP 表示发送挂起（SIGHUP）信号，通常用于通知守护进程重新加载配置文件或重启自身。
  - $MAINPID 是 systemd 提供的一个特殊变量，代表服务的主进程 ID（PID）。这确保了信号被发送到正确的进程。
* WatchdogSec 启用 Watchdog，systemd 将在 10 秒内未收到服务的心跳（通知）时，自动重启服务。
  systemd 将在 10 秒后认为服务可能已经挂起或卡死，如果在这 10 秒内服务没有通过特定的机制（如调用 sd_notify("WATCHDOG=1")）来“喂狗”（即重置计时器），systemd 将自动重启服务。
  启用 Watchdog 的前提条件：
  - 服务必须以 Type=notify 启动，这样 systemd 才能接收来自服务的通知。
  - 服务本身需要定期调用 sd_notify("WATCHDOG=1") 来重置 Watchdog 计时器。如果服务未实现这一机制，Watchdog 功能将无法正常工作，可能导致服务被误杀，将导致服务被频繁重启。
  - 调试: 在开发和测试阶段，可以暂时减少 WatchdogSec 的值（如 5 秒）以便更快地观察行为，但在生产环境中应使用适当的值。
* NotifyAccess=all 允许服务发送所有类型的通知给 systemd，包括 Watchdog 通知。如果你的服务仅需要发送 Watchdog 通知，可以将其设置为 NotifyAccess=watchdog 以增强安全性。

Install
* WantedBy=multi-user.target 确保服务在多用户运行级别下启动。


```sh
# 重新加载 systemd 配置
sudo systemctl daemon-reload
# 启动服务
sudo systemctl start your-service.service
# 检查服务状态
sudo systemctl status your-service.service
# 重新加载服务
sudo systemctl reload your-service.service
# 查看日志
journalctl -u your-service.service -f
```

# 监控指标
```go
// Prometheus监控埋点
var(
    restartsTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "endless_restarts_total",
        Help: "Total number of graceful restarts",
    })
    activeConnections = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "endless_active_connections",
        Help: "Currently active connections",
    })
)
```
