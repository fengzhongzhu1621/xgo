package tpl_test

import (
	"os"
	"testing"
	"text/template"
)

func TestRange(t *testing.T) {
	const templateText = `
<ul>
{{range .Projects}}<li>{{.}}</li>{{end}}
</ul>
`
	tmpl, err := template.New("TestRange").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Projects []string
	}{
		Projects: []string{"foo", "bar"},
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func TestRangeMap(t *testing.T) {
	const templateText = `
	<ul>
	{{range $key, $value := .Projects}}<li>{{$key}}: {{$value}}</li>{{end}}
	</ul>
	`
	tmpl, err := template.New("TestRangeMap").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Projects map[string]string
	}{
		Projects: map[string]string{"Name": "bob", "Address": "address"},
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
