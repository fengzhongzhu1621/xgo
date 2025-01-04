package tpl

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestTemplateReplace 用数据映射中的相应值替换模板字符串中的占位符。占位符用大括号括起来，例如 {key}。
// 例如，模板字符串是“Hello, {name}!”，数据映射是{"name": "world"}，结果将是“Hello, world!”。
// func TemplateReplace(template string, data map[string]string string
func TestTemplateReplace(t *testing.T) {
	template := `Hello, my name is {name}, I'm {age} years old.`
	data := map[string]string{
		"name": "Bob",
		"age":  "20",
	}

	result := strutil.TemplateReplace(template, data)

	assert.Equal(t, "Hello, my name is Bob, I'm 20 years old.", result, "result")
}
