package metrics

// 标签（Tags）设计
//
// 调用关系标识
const (
	// TagCaller is the caller name.
	TagCaller = "caller"
	// TagCallee is the callee name.
	TagCallee = "callee"
	// TagMethod is the callee method.
	// 记录调用的具体方法（如API接口名），细化监控粒度
	TagMethod = "method"
)

// 重试与流量控制
const (
	// TagAttempts gives how many attempts an app request has used.
	// 记录请求的重试次数，结合 FQNRealRequest 可分析重试策略有效性。
	TagAttempts = "attempts"
	// TagThrottled indicate whether is this request throttled.
	// 标记请求是否被限流（如令牌桶拒绝），反映系统过载保护状态
	TagThrottled = "throttled"
	// TagNoMoreAttempt indicate whether the server issues no retry/hedging.
	// 标识服务端禁止重试/对冲，用于调试异常场景
	TagNoMoreAttempt = "noMoreAttempt"
)

// 错误与状态
const (
	// TagErrCodes is the error code the request.
	// 记录请求的错误码（如HTTP 500或自定义错误），分类统计故障原因
	TagErrCodes = "error_codes"
	// TagInflight is inflight number of app request or indicate whether the real request is still inflight.
	// 实时统计正在处理的请求数（应用层或真实请求），监控系统并发压力。
	TagInflight = "inflight"
)

// TagNamesApp gives the order in which tagPairs are organized for FQNAppRequest or FQNAppCostMs.
// 为应用层指标（FQNAppRequest 和 FQNAppCostMs）定义标签的组织顺序，反映用户视角的请求行为。
var TagNamesApp = []string{
	TagCaller,
	TagCallee,
	TagMethod,
	TagAttempts,      // 记录重试次数，用于分析重试策略有效性。
	TagErrCodes,      // 分类错误类型，辅助故障诊断。
	TagThrottled,     // 标记是否触发限流，反映系统负载状态。
	TagInflight,      // 实时统计并发请求数。
	TagNoMoreAttempt, // 标识服务端禁止重试，用于调试异常场景。
}

// TagNamesReal gives the order in which tagPairs are organized for FQNRealRequest or FQNRealCostMs.
// 为真实请求指标（FQNRealRequest 和 FQNRealCostMs）定义标签顺序，聚焦于实际调用细节。
var TagNamesReal = []string{
	TagCaller,
	TagCallee,
	TagMethod,
	TagErrCodes,
	TagInflight,
	TagNoMoreAttempt,
}
