package tpl

import (
	"html/template"
	"net/http"

	"github.com/fengzhongzhu1621/xgo"
)

var _tmpls = make(map[string]*template.Template)

// RenderHTML 根据文件名，读取模板文件，并使用 v 进行渲染模板内容
func RenderHTML(w http.ResponseWriter, name string, v interface{}) {
	if t, ok := _tmpls[name]; ok {
		// 使用模板渲染
		t.Execute(w, v)
		return
	}
	// 读取模板文件，渲染资源文件的内容
	t := template.Must(
		template.New(name).Funcs(xgo.FuncMap).Delims("[[", "]]").Parse(xgo.ReadAssetsContent(name)),
	)
	// 添加到缓存中
	_tmpls[name] = t
	// 使用上下文变量 v 渲染资源文件的内容，将渲染的内容写入到 http 响应中
	t.Execute(w, v)
}
