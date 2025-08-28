package debuglog

// RuleItem is the basic configuration for a single rule.
type RuleItem struct {
	Method  *string // 目标方法名（如 "GET"、"POST"），nil 表示不限制; 指向字符串的指针，允许nil（表示忽略该字段）
	Retcode *int    // 目标返回码（如 200、404），nil 表示不限制; 指向整数的指针，允许nil（表示忽略该字段）
}

// Matched is the result of rule matching.
func (e RuleItem) Matched(destMethod string, destRetCode int) bool {
	// 如果 Method 或 Retcode 为 nil，则跳过对应字段的匹配。
	// 否则，严格比较值是否相等
	// e.g. /trpc.app.server.service/methodA
	return (e.Method == nil || *e.Method == destMethod) &&
		(e.Retcode == nil || *e.Retcode == destRetCode)
}
