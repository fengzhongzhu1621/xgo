package json

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// omitempty：这个选项告诉编码器，如果该字段的值是零值（对于数值类型是0，字符串类型是空字符串，布尔类型是false等），
// 则在编码为JSON时忽略该字段。换句话说，如果字段值为空或零值，它不会出现在生成的JSON对象中。
// 这在减少JSON数据的大小和提高传输效率方面很有用。

// omitempty 无法处理空 struct，例如 Post.Category。
// omitempty 处理 time.Time 的方式并非我们理解的 UTC = 0，即 1970-01-01 00:00:00，而是 0001-01-01T00:00:00Z。

type Post struct {
	Id         int64           `json:"id,omitempty"`
	CreateTime time.Time       `json:"create_time,omitempty"` // 0001-01-01T00:00:00Z
	TagList    []Tag           `json:"tag_list,omitempty"`
	Name       string          `json:"name,omitempty"`
	Score      float64         `json:"score,omitempty"`
	Category   Category        `json:"category,omitempty"` // {"id":0,"name":""}
	LikePost   map[string]bool `json:"like,omitempty"`     // 改为记录用户是否点赞
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int64  `json:"id"` // 改为int64
	Name string `json:"name"`
}

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

func TestOmitmptyTimeAndStruct(t *testing.T) {
	// 序列化为JSON
	b, err := json.Marshal(new(Post))
	if err != nil {
		fmt.Println("JSON序列化失败:", err)
		return
	}

	// 打印结果
	// {"create_time":"0001-01-01T00:00:00Z","category":{"id":0,"name":""}}
	fmt.Println(string(b))
}
