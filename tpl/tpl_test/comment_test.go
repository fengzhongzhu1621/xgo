package tpl_test

import (
	"os"
	"testing"
	"text/template"
)

// 注释内容不会出现在最终输出中，可以用于在模板中添加说明或备注。
func TestComment(t *testing.T) {
	const templateText = `Hello, {{.Name}}! {{/* This is a comment */}}`
	tmpl, err := template.New("TestComment").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string
	}{
		Name: "bob",
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
