package xgo

import (
	"embed"
	"net/http"
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