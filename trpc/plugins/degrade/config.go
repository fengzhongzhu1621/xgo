package degrade

var cfg Config

// Config 熔断配置结构体声明
type Config struct {
	// Load5 5 分钟 触发熔断的阈值，恢复时使用实时 load1 <= 本值来判断，更敏感
	Load5 float64 `yaml:"load5"`
	// CPUIdle cpuidle 触发熔断的阈值
	CPUIdle int `yaml:"cpu_idle"`
	// MemoryUsepercent 内存使用率百分比 触发熔断的阈值
	MemoryUsePercent int `yaml:"memory_use_p"`
	// DegradeRate 流量保留比例，目前使用随机算法抛弃，迭代加入其他均衡算法；0 或 100 则不启用此插件; -1 表示丢弃全部流量
	DegradeRate int `yaml:"degrade_rate"`
	// Interval 心跳时间间隔，主要控制多久更新一次熔断开关状态
	Interval int `yaml:"interval"`
	// Modulename 模块名，后续用于上报鹰眼或其他日志
	Modulename string `yaml:"modulename"`
	// Whitelist 白名单，用于后续跳过不被熔断控制的业务接口
	Whitelist string `yaml:"whitelist"`
	// IsActive 标志熔断是否生效
	IsActive bool `yaml:"-"`
	// MaxConcurrentCnt 最大并发请求数，<=0 时不开启，控制最大并发请求数，和上述熔断互为补充，能防止突发流量把服务打死，比如 1ms 内突然进入 100W 请求
	MaxConcurrentCnt int `yaml:"max_concurrent_cnt"`
	// MaxTimeOutMs 超过最大并发请求数时，最多等待 MaxTimeOutMs 才决定是熔断还是继续处理
	MaxTimeOutMs int `yaml:"max_timeout_ms"`
}

// enableConcurrency 是否开启最大并发请求数的限制
func enableConcurrency() bool {
	return cfg.MaxConcurrentCnt > 0
}

// GetConfig 获取配置参数
func GetConfig() Config {
	return cfg
}
