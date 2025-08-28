package debuglog

type options struct {
	logFunc         LogFunc      // 日志记录函数
	errLogLevelFunc LogLevelFunc // 错误日志级别判断函数
	nilLogLevelFunc LogLevelFunc // 空值日志级别判断函数
	enableColor     bool         // 是否启用颜色输出
	include         []*RuleItem  // 包含规则列表
	exclude         []*RuleItem  // 排除规则列表
}

// passed is the filtering result.
// When it is true, will go to the logging process.
func (o *options) passed(rpcName string, errCode int) bool {
	// Calculation of the include rule.
	// 包含规则优先
	for _, in := range o.include {
		if in.Matched(rpcName, errCode) {
			return true
		}
	}
	// If the include rule is configured, the exclude rule will not be matched.
	// 明确包含规则存在但未匹配时拒绝
	if len(o.include) > 0 {
		return false
	}

	// 没有包含规则时，匹配排除规则
	// Calculation of the exclude rule.
	for _, ex := range o.exclude {
		if ex.Matched(rpcName, errCode) {
			return false
		}
	}

	return true
}

// getFilterOptions gets the interceptor condition options.
func getFilterOptions(opts ...Option) *options {
	o := &options{
		logFunc:         DefaultLogFunc,               // 默认打印请求和响应结构体的值
		errLogLevelFunc: LogContextfFuncs[errorLevel], // 默认日志级别为 error
		nilLogLevelFunc: LogContextfFuncs[debugLevel], // 默认日志级别为 debug
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
