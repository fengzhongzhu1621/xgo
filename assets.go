package xgo

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

// embed.FS 是 Go 1.16 版本引入的一个新功能，它是 embed 包中的一个类型。embed.FS 提供了一种将文件或文件夹嵌入到 Go 二进制程序中的方法。这意味着你可以在编译时将静态资源（如 HTML、CSS、JavaScript 文件或图像）打包到你的 Go 应用程序中，
// 并在运行时直接访问它们，而无需从外部文件系统或网络加载。

//go:embed assets
var assetsFS embed.FS

// http.FS 是 Go 1.16 中引入的另一个新功能，位于 net/http 包中。
// http.FS 是一个基于 HTTP/HTTPS 协议的文件系统接口，它允许你通过 HTTP/HTTPS URL 访问文件系统中的文件和目录。
// http.FS 提供了一种灵活的方式来访问远程文件系统，而无需显式地下载文件到本地。你可以使用 http
// Assets contains project assets.
var Assets = http.FS(assetsFS)


// 读取资源文件的内容
func ReadAssetsContent(name string) string {
	// 打开资源文件
	fd, err := Assets.Open(name)
	if err != nil {
		panic(err)
	}
	// 读取资源文件的内容
	data, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	// 将字节数组转换为字符串
	return string(data)
}


var (
	FuncMap template.FuncMap // 定义的模板函数
)


func init() {
	FuncMap = template.FuncMap{
		"title": strings.ToTitle, // 首字母大写
		"urlhash": func(path string) string {
			// 打开资源文件
			httpFile, err := Assets.Open(path)
			if err != nil {
				return path + "#no-such-file"
			}
			// 构造文件的 uri
			info, err := httpFile.Stat()
			if err != nil {
				return path + "#stat-error"
			}
			return fmt.Sprintf("%s?t=%d", path, info.ModTime().Unix())
		},
	}
}
