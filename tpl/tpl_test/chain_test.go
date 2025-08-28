package tpl_test

import (
	"os"
	"testing"
	"text/template"
)

func TestChain(t *testing.T) {
	const templateText = `Name: {{.Repo.Name}}, Address: {{.Repo.Address}}`
	tmpl, err := template.New("TestChain").Parse(templateText)
	if err != nil {
		panic(err)
	}

	data := struct {
		Repo struct {
			Name    string
			Address string
		}
	}{
		Repo: struct {
			Name    string
			Address string
		}{
			Name:    "bob",
			Address: "address",
		},
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
