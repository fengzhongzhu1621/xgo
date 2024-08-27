package tpl_test

import (
	"os"
	"testing"
	"text/template"
)

func TestDefine(t *testing.T) {
	const (
		headerTemplate = `{{define "header"}}<html><head><title>{{.Title}}</title></head><body>{{end}}`
		footerTemplate = `{{define "footer"}}</body></html>{{end}}`
		bodyTemplate   = `{{define "body"}}<h1>{{.Heading}}</h1><p>{{.Content}}</p>{{end}}`
		mainTemplate   = `{{template "header" .}} {{template "body" .}} {{template "footer" .}}`
	)

	tmpl := template.Must(template.New("TestDefine").Parse(headerTemplate + footerTemplate + bodyTemplate + mainTemplate))

	data := struct {
		Title   string
		Heading string
		Content string
	}{
		Title:   "Welcome",
		Heading: "Hello, bob!",
		Content: "This is a simple example of nested templates.",
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
