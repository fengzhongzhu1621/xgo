package tpl_test

import (
	"os"
	"testing"
	"text/template"
)

func TestIf(t *testing.T) {
	const templateText = `{{if .ShowTitle}}<h1>{{.Title}}</h1>{{end}}`
	tmpl, err := template.New("TestIf").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		ShowTitle bool
		Title     string
	}{
		ShowTitle: true,
		Title:     "Hello, bob!",
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func TestElse(t *testing.T) {
	const templateText = `{{if .ShowTitle}}<h1>{{.Title}}</h1>{{else}}<h1>No Title</h1>{{end}}`
	tmpl, err := template.New("TestElse").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		ShowTitle bool
		Title     string
	}{
		ShowTitle: false,
		Title:     "Hello, bob!",
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func TestElseIf(t *testing.T) {
	const templateText = `{{if .ShowTitle}}<h1>{{.Title}}</h1>{{else if .ShowSubtitle}}<h2>{{.Subtitle}}</h2>{{else}}<p>No Title or Subtitle</p>{{end}}`
	tmpl, err := template.New("example").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		ShowTitle    bool
		ShowSubtitle bool
		Title        string
		Subtitle     string
	}{
		ShowTitle:    false,
		ShowSubtitle: true,
		Title:        "Hello, bob!",
		Subtitle:     "Hi, bob!",
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
