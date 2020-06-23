package xgo

import "os"

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

/**
 * 从命令行获得服务器的IP和端口
 */
func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

/**
 * 返回整数的指针
 */
func IntPtr(i int)  *int {
	return &i
}
