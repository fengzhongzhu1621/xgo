package degrade

import (
	"time"

	"trpc.group/trpc-go/tnet/log"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginType = "circuitbreaker"
	pluginName = "degrade"
)

const (
	infoDegradeRateZero = "because the derade_rate is zero,so exit plugin"
	infoDegradeRate100  = "the derade_rate is 100, so exit plugin"
	errDegardeReturn    = "service is degrade..."
	systemDegradeErrNo  = 22
)

var (
	isDegrade bool
	sema      chan struct{}
)

// Degrade 熔断插件默认初始化
type Degrade struct{}

// Type 返回插件类型
func (p *Degrade) Type() string {
	return pluginType
}

// Setup 注册
func (p *Degrade) Setup(name string, decoder plugin.Decoder) error {
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	// 流量抛弃比例，目前使用随机算法抛弃，迭代加入其他均衡算法
	if cfg.DegradeRate == 0 {
		log.Info(infoDegradeRateZero)
		return nil
	}
	if cfg.DegradeRate == 100 {
		log.Info(infoDegradeRate100)
		return nil
	}

	// 心跳时间间隔，主要控制多久更新一次熔断开关状态
	if cfg.Interval == 0 {
		cfg.Interval = 60
	}

	// 设置最大并发请求数
	if enableConcurrency() {
		sema = make(chan struct{}, cfg.MaxConcurrentCnt)
	}

	// 异步更新系统数据给全局变量，目前只更新 cpu 空闲率
	go UpdateSysInfoPerTime()

	// 异步更新熔断状态
	go func() {
		var load1, load5 float64
		for range time.Tick(time.Duration(cfg.Interval) * time.Second) {
			// 周期检测获取系统负载情况
			cpuIdle := GetCPUIdle()
			mem := int(GetMemoryStat())
			load, err := GetLoadAvg()

			// 根据 CPUidle，内存使用率，负载（主要 load5）来设置阈值，达到阈值触发熔断保护，抛弃一定百分比的随机流量
			if err != nil {
				load5 = 0
				load1 = 0
			} else {
				load1 = load.Load1
				load5 = load.Load5
			}
			if load5 > cfg.Load5 || mem > cfg.MemoryUsePercent || cpuIdle < cfg.CPUIdle {
				isDegrade = true
			}
			// 恢复时使用实时 load1 <= 本值来判断，更敏感
			if load1 <= cfg.Load5 && mem <= cfg.MemoryUsePercent && cpuIdle >= cfg.CPUIdle {
				isDegrade = false
			}
			log.Infof("%s cpu_idle:%d mem_usage:%d load5:%f,degrade:%t",
				time.Now(), cpuIdle, mem, load5, isDegrade)
		}
	}()

	filter.Register(pluginName, Filter, nil)

	return nil
}
