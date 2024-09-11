package xgo

import (
	"os"
	"reflect"
)

// 空值定义.
var (
	EmptyStr    string
	EmptyError  error
	EmptyResult []interface{}
	EmptyArgs   []reflect.Value
)

// ResolveAddress 从命令行获得服务器的IP和端口.
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

func Assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}
