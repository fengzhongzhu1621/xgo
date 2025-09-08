package xerror

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// frame 表示堆栈帧内的一个程序计数器。
// 由于历史原因，如果将 frame 解释为 uintptr，其值表示程序计数器 + 1。
type frame uintptr

// pc 返回此帧的程序计数器；
// 多个帧可能具有相同的 PC 值。
func (f frame) pc() uintptr {
	return uintptr(f) - 1 // 计算实际的程序计数器值（减去1）
}

// file 返回包含该帧 pc 对应函数的文件的完整路径。
func (f frame) file() string {
	// 根据程序计数器（PC）地址返回对应的函数信息
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	// 根据程序计数器（PC）地址返回对应的文件信息
	file, _ := fn.FileLine(f.pc())
	return file
}

// line 返回该帧 pc 对应函数的源代码行号。
func (f frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	// 根据程序计数器（PC）地址返回对应的行号信息
	_, line := fn.FileLine(f.pc())
	return line
}

// name 返回此函数的名称（如果已知）。
func (f frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	// 根据程序计数器（PC）地址返回对应的函数名称
	return fn.Name()
}

// Format 根据 fmt.Formatter 接口格式化帧。
//
//	%s    源文件
//	%d    源代码行号
//	%n    函数名
//	%v    等同于 %s:%d
//
// Format 接受改变某些动词输出的标志：
//
//	%+s   函数名和相对于编译时 GOPATH 的源文件路径，以 \n\t 分隔 (<funcName>\n\t<path>)
//	%+v   等同于 %+s:%d
func (f frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's': // 处理文件路径输出
		switch {
		case s.Flag('+'): // 详细模式：输出函数名和完整文件路径
			io.WriteString(s, f.name()) // 写入函数名
			io.WriteString(s, "\n\t")   // 写入换行和制表符
			io.WriteString(s, f.file()) // 写入文件完整路径
		default: // 默认模式：只输出基础文件名
			io.WriteString(s, path.Base(f.file())) // 使用 path.Base 获取路径最后一部分
		}
	case 'd': // 处理行号输出
		io.WriteString(s, strconv.Itoa(f.line())) // 将行号(int)转换为字符串写入
	case 'n': // 处理函数名输出
		io.WriteString(s, funcName(f.name())) // 写入处理后的函数名（去除路径和包名）
	case 'v': // 处理默认值输出
		f.Format(s, 's')       // 先格式化文件部分
		io.WriteString(s, ":") // 写入冒号分隔符
		f.Format(s, 'd')       // 再格式化行号部分
	}
}

// funcName 移除由 func.Name() 报告的函数名的路径前缀组件。
func funcName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
