package httpstaticserver

import (
	"bytes"
	"html/template"
	"runtime"
)

// versionMessage 获得静态服务的版本信息
func versionMessage() string {
	// 创建字符串模板
	t := template.Must(template.New("version").Parse(`GoHTTPServer
  Version:        {{.Version}}
  Go version:     {{.GoVersion}}
  OS/Arch:        {{.OSArch}}
  Git commit:     {{.GitCommit}}
  Built:          {{.Built}}
  Site:           {{.Site}}`))
	// 渲染模板，将结果输出到 buf 中
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, map[string]interface{}{
		"Version":   VERSION,
		"GoVersion": runtime.Version(),
		"OSArch":    runtime.GOOS + "/" + runtime.GOARCH,
		"GitCommit": GITCOMMIT,
		"Built":     BUILDTIME,
		"Site":      SITE,
	})

	// 将 bytes.Buffer 转换为 string
	return buf.String()
}
