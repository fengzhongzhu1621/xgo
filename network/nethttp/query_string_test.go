package nethttp

import (
	"fmt"
	"testing"

	"github.com/google/go-querystring/query"
	goutil_maputil "github.com/gookit/goutil/maputil"
	"github.com/stretchr/testify/assert"
)

// 字典转换为querystring
func TestHTTPQueryString(t *testing.T) {
	src := map[string]any{"a": "v0", "b": 23}
	str := goutil_maputil.HTTPQueryString(src)

	fmt.Println(str) // a=v0&b=23

	assert.Contains(t, str, "b=23")
	assert.Contains(t, str, "a=v0")
}

// 结构体转换为querystring
func TestQuerystringEncode(t *testing.T) {
	// 注意：数据结构属性名需要大写
	type Data struct {
		Name      string `url:"name"`
		Age       int    `url:"age"`
		IsStudent bool   `url:"isStudent"`
	}

	data := Data{
		Name:      "Tom",
		Age:       2,
		IsStudent: true,
	}

	value, _ := query.Values(data)

	output := value.Encode()
	fmt.Println(output) // age=2&isStudent=true&name=Tom
}
