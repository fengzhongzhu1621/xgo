package tpl

import "html/template"

var (
	tmpl *template.Template
	// 定义了模板缓存
	templates = map[string]string{
		"index":       "res/index.tmpl.html",
		"ipa-install": "res/ipa-install.tmpl.html",
	}
)

// 创建一个名为 name 的模板对象，解析模板内容 content
func ParseTemplate(name string, content string) {
	if tmpl == nil {
		// 创建一个新的模板对象
		tmpl = template.New(name)
	}
	var t *template.Template
	if tmpl.Name() == name {
		t = tmpl
	} else {
		t = tmpl.New(name)
	}
	// 自定义模板分隔符，使用 Parse() 解析模板内容，使用 Must 函数处理潜在的错误
	template.Must(t.New(name).Delims("[[", "]]").Parse(content))
}
