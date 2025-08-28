package tpl_test

import (
	"os"
	"strings"
	"testing"
	"text/template"
)

func TestFuncMap(t *testing.T) {
	const templateText = `
<p>Upper: {{.Name | upper}}</p>
<p>Len: {{len .Name}}</p>
`

	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}

	tmpl, err := template.New("TestFuncMap").Funcs(funcMap).Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string
	}{
		Name: "John",
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}

func TestRepeat(t *testing.T) {
	funcMap := template.FuncMap{
		"repeat": func(s string, count int) string {
			return strings.Repeat(s, count)
		},
	}

	const templateText = `{{repeat .Name 3}}`
	tmpl, err := template.New("TestRepeat").Funcs(funcMap).Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string
	}{
		Name: "Go",
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}

func TestPipeline(t *testing.T) {
	funcMap := template.FuncMap{
		"trim":  strings.TrimSpace,
		"upper": strings.ToUpper,
	}

	const templateText = `{{.Name | trim | upper}}`
	tmpl, err := template.New("TestPipeline").Funcs(funcMap).Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string
	}{
		Name: "  go  ",
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
