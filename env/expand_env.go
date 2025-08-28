package env

import (
	"os"
)

// ExpandEnv looks for ${var} in s and replaces them with value of the corresponding environment variable.
// It's not like os.ExpandEnv which handles both ${var} and $var.
// $var is considered invalid, since configurations like password for redis/mysql may contain $.
// ExpandEnv 在字节切片中查找 ${var} 模式，并用对应的环境变量值替换
// 与标准库 os.ExpandEnv 不同，这里只处理 ${var} 格式，不处理 $var 格式
// 因为像redis/mysql密码等配置可能包含$字符，$var格式容易产生误匹配
func ExpandEnv(s []byte) []byte {
	var buf []byte // 用于构建结果的可变字节切片
	i := 0         // 当前处理位置的起始索引
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+2 < len(s) && s[j+1] == '{' { // 只识别 ${var} 格式，不识别 $var 格式
			if buf == nil {
				buf = make([]byte, 0, 2*len(s)) // 预分配足够容量的缓冲区
			}
			buf = append(buf, s[i:j]...)   // 添加${之前的普通文本
			name, w := getEnvName(s[j+1:]) // 获取环境变量名和匹配长度
			if name == nil && w > 0 {
				// 无效匹配，移除$符号
			} else if name == nil {
				buf = append(buf, s[j]) // 保留$符号
			} else {
				buf = append(buf, os.Getenv(string(name))...) // 用环境变量值替换${var}
			}
			j += w    // 跳过已处理的环境变量部分
			i = j + 1 // 更新起始索引
		}
	}
	if buf == nil {
		return s // 如果没有环境变量替换，返回原切片
	}
	return append(buf, s[i:]...) // 添加剩余文本并返回
}

// getEnvName gets env name, that is, var from ${var}.
// The env name and its len will be returned.
// getEnvName 从 ${var} 格式中提取环境变量名
// 返回环境变量名和匹配的总长度
func getEnvName(s []byte) ([]byte, int) {
	// 查找右花括号 '}'
	// 保证第一个字符是 '{' 且字符串至少有两个字符
	for i := 1; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\n' || s[i] == '"' { // "xx${xxx" 遇到无效字符
			return nil, 0 // 遇到无效字符，保留$符号
		}
		if s[i] == '}' {
			if i == 1 { // ${} 空的花括号
				return nil, 2 // 移除 ${}
			}
			return s[1:i], i + 1 // 返回环境变量名和总长度（包括${和}）
		}
	}
	return nil, 0 // 没有找到}，保留$符号
}
