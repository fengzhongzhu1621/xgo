package xgo

import (
	"os"
	"reflect"
)

// 空值定义
var (
	EmptyStr    string
	EmptyError  error
	EmptyResult []interface{}
	EmptyArgs   []reflect.Value
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

// 从命令行获得服务器的IP和端口
func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrintf("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrintf("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

// 返回整数的指针
func IntPtr(i int) *int {
	return &i
}

func Assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func LastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
