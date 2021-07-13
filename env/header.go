package env

import (
	"fmt"
	"strings"
	"xgo/utils/string_utils"
)

/**
 * 构造header
 * 返回一个新数组
 */
func AppendEnv(env []string, k string, v ...string) []string {
	if len(v) == 0 {
		return env
	}
	// 创建一个字符串空数组
	vCleaned := make([]string, 0, len(v))
	// 将数组元素去掉换行符和首尾的空白字符
	for _, val := range v {
		vCleaned = append(vCleaned, strings.TrimSpace(string_utils.HeaderNewlineToSpace.Replace(val)))
	}
	return append(env, fmt.Sprintf("%s=%s",
		strings.ToUpper(k),
		strings.Join(vCleaned, ", ")))
}
