package httpStaticServer

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// 获得静态路由执行的文件所在的绝对路径
func (s *HTTPStaticServer) getRealPath(r *http.Request) string {
	// 获得请求路由的路径
	path := mux.Vars(r)["path"]
	// 给路由增加前缀 /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// 用于清理文件路径。它会去除路径中多余的斜杠、解析 . 和 .. 等特殊序列，并返回规范化后的绝对路径。
	path = filepath.Clean(path) // prevent .. for safe issues
	// 用于计算两个文件路径之间的相对路径
	relativePath, err := filepath.Rel(s.Prefix, path)
	if err != nil {
		relativePath = path
	}
	// 计算文件所在的绝对路径
	realPath := filepath.Join(s.Root, relativePath)

	// 将文件路径中的操作系统特定的分隔符转换为正斜杠 /
	return filepath.ToSlash(realPath)
}