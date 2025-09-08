package xerror

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
	traceable bool               // 控制是否启用堆栈跟踪功能
	content   string             // 用于过滤堆栈信息的内容字符串
	stackSkip = defaultStackSkip // 收集堆栈时跳过的帧数
)

const (
	defaultStackSkip = 3 // 默认跳过3层堆栈帧（通常跳过xerror包自身的调用及runtime包调用）
)

// SetTraceable 控制是否在错误中启用堆栈跟踪
func SetTraceable(x bool) {
	traceable = x
}

// SetTraceableWithContent 启用堆栈跟踪并设置过滤内容
// 打印堆栈信息时，只输出包含指定内容的内容
// 可通过将content设置为服务名来过滤其他插件的堆栈信息，避免输出大量无用信息[2](@ref)
func SetTraceableWithContent(c string) {
	traceable = true // 启用堆栈跟踪
	content = c      // 设置过滤内容
}

// SetStackSkip 设置收集堆栈时要跳过的帧数
// 在封装New方法时，可根据封装层数设置stackSkip（例如设置为4）
// 此函数在项目启动前设置，不保证并发安全
func SetStackSkip(skip int) {
	stackSkip = skip
}

// stackTrace 表示堆栈跟踪，从最内层（最新）到最外层（最旧）的帧集合
type stackTrace []frame

// Format 根据fmt.Formatter接口格式化堆栈帧
//
//	%s    列出堆栈中每个帧的源文件
//	%v    列出堆栈中每个帧的源文件和行号
//
// Format 接受改变某些动词输出的标志：
//
//	%+v   打印堆栈中每个帧的文件名、函数名和行号
func (st stackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'): // 详细模式
			for _, f := range st {
				// 过滤：只打印包含指定内容的堆栈信息
				if !isOutput(fmt.Sprintf("%+v", f)) {
					continue // 如果不包含过滤内容，则跳过该帧
				}
				io.WriteString(s, "\n") // 每个帧前换行
				f.Format(s, verb)       // 格式化单个帧
			}
		case s.Flag('#'): // Go语法表示模式
			fmt.Fprintf(s, "%#v", []frame(st)) // 输出堆栈切片的Go语法表示
		default: // 默认模式
			st.formatSlice(s, verb)
		}
	case 's': // 字符串模式
		st.formatSlice(s, verb)
	}
}

// formatSlice 将stackTrace格式化为帧切片形式
// 仅在'%s'或'%v'动词时有效
func (st stackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[") // 开始切片
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ") // 帧之间用空格分隔
		}
		f.Format(s, verb) // 格式化单个帧
	}
	io.WriteString(s, "]") // 结束切片
}

// isOutput 检查字符串是否包含指定的内容，用于过滤堆栈信息
func isOutput(str string) bool {
	return strings.Contains(str, content)
}

// callers 获取当前调用堆栈
func callers() stackTrace {
	const depth = 32                        // 最大堆栈深度
	var pcs [depth]uintptr                  // 存储程序计数器（PC）的数组
	n := runtime.Callers(stackSkip, pcs[:]) // 跳过stackSkip层获取堆栈
	stack := pcs[0:n]                       // 截取有效部分

	// 将uintptr数组转换为frame切片
	st := make([]frame, len(stack))
	for i := 0; i < len(st); i++ {
		st[i] = frame((stack)[i]) // 转换每个程序计数器为frame
	}

	return st // 返回堆栈跟踪
}
