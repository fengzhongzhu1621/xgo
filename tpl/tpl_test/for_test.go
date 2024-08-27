package tpl_test

import (
	"os"
	"testing"
	"text/template"
)

func TestBreak(t *testing.T) {
	const templateText = `
	<ul>{{range .Items}}
	  {{if eq . "Item 3"}}{{break}}{{end}}<li>{{.}}</li>{{end}}
	</ul>`

	tmpl, err := template.New("TestBreak").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Items []string
	}{
		Items: []string{"Item 1", "Item 2", "Item 3"},
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func TestContinue(t *testing.T) {
	const templateText = `<ul>{{range .Items}}
	{{if eq . "Item 2"}}{{continue}}{{else}}<li>{{.}}</li>{{end}}{{end}}
  </ul>`

	tmpl, err := template.New("TestContinue").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Items []string
	}{
		Items: []string{"Item 1", "Item 2", "Item 3"},
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
