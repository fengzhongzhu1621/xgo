https://github.com/go-gorm/prometheus
https://github.com/go-gorm/prometheus/blob/master/mysql.go

```go
import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "gorm.io/plugin/prometheus"
)

db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

db.Use(prometheus.New(prometheus.Config{
  DBName:          "db1", // 使用 `DBName` 作为指标 label
  RefreshInterval: 15,    // 指标刷新频率（默认为 15 秒）
  PushAddr:        "prometheus pusher address", // 如果配置了 `PushAddr`，则推送指标
  StartServer:     true,  // 启用一个 http 服务来暴露指标
  HTTPServerPort:  8080,  // 配置 http 服务监听端口，默认端口为 8080 （如果您配置了多个，只有第一个 `HTTPServerPort` 会被使用）
  MetricsCollector: []prometheus.MetricsCollector {
    &prometheus.MySQL{
      VariableNames: []string{"Threads_running"},
    },
  },  // 用户自定义指标
}))
```
