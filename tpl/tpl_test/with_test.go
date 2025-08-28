package tpl_test

import (
	"os"
	"testing"
	"text/template"
)

// with action 设置一个新的数据上下文，并在该上下文中执行模板。
// with action 将 User 结构体设置为新的数据上下文{{.}}，从而简化了模板中对嵌套字段的访问。
func TestWith(t *testing.T) {
	const templateText = `{{with .User}}<p>Name: {{.Name}}</p>
<p>Age: {{.Age}}</p>
{{end}}
`

	tmpl, err := template.New("TestWith").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		User struct {
			Name string
			Age  int
		}
	}{
		User: struct {
			Name string
			Age  int
		}{
			Name: "bob",
			Age:  30,
		},
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
