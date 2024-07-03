package tpl

import (
	"html/template"
	"net/http"

	"github.com/fengzhongzhu1621/xgo"
)

var (
	_tmpls = make(map[string]*template.Template)
)


// 根据文件名，读取模板文件，并使用 v 进行渲染模板内容
func RenderHTML(w http.ResponseWriter, name string, v interface{}) {
	if t, ok := _tmpls[name]; ok {
		// 使用模板渲染
		t.Execute(w, v)
		return
	}
	// 读取模板文件
	t := template.Must(template.New(name).Funcs(xgo.FuncMap).Delims("[[", "]]").Parse(xgo.ReadAssetsContent(name)))
	_tmpls[name] = t
	t.Execute(w, v)
}
