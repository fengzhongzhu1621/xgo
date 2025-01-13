package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

// omitempty：这个选项告诉编码器，如果该字段的值是零值（对于数值类型是0，字符串类型是空字符串，布尔类型是false等），
// 则在编码为JSON时忽略该字段。换句话说，如果字段值为空或零值，它不会出现在生成的JSON对象中。
// 这在减少JSON数据的大小和提高传输效率方面很有用。

type Example struct {
	Data *string `json:"data,omitempty"`
}

func TestOmitmpty(t *testing.T) {
	example := Example{}

	// 当Data为nil时，编码后的JSON不包含"data"字段
	jsonData, _ := json.Marshal(example)
	fmt.Println(string(jsonData)) // 输出：{}

	str := "Hello, World!"
	example.Data = &str

	// 当Data不为nil时，编码后的JSON包含"data"字段
	jsonData, _ = json.Marshal(example)
	fmt.Println(string(jsonData)) // 输出：{"data":"Hello, World!"}
}
