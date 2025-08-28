package metrics

// 请求量统计
const (
	// FQNAppRequest the number of requests triggered by user.
	// 记录用户触发的请求数（应用层视角），用于衡量用户行为频率或系统入口流量
	FQNAppRequest = "appRequest"
	// FQNRealRequest the number of total requests sent to callee.
	// 实际发送给被调用方（如微服务、数据库）的请求总数，包含重试/对冲请求，反映真实负载
	FQNRealRequest = "realRequest"
)

// 耗时统计
const (
	// FQNAppCostMs the cost of app requests.
	// 应用层请求总耗时（从用户触发到最终响应），用于衡量用户体验
	FQNAppCostMs = "appCostMs"
	// FQNRealCostMs the cost of real requests.
	// 单个真实请求的耗时（如某次重试的RPC调用时间），帮助定位性能瓶颈
	FQNRealCostMs = "realCostMs"
)
