package gjson

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

// GetMany()可以一次读取多个值，返回一个数组[]gjson.Result
// gjson 在解析的过程中也不会像 fastjson 一样将解析的内容保存在一个结构体中，可以反复的利用。
// 所以当调用 GetMany 想要返回多个值的时候，其实也是需要遍历 JSON 串多次，因此效率会比较低。
func TestGetMany(t *testing.T) {
	results := gjson.GetMany(json4, "name", "age", "pets.#", "contact.phone")
	for _, result := range results {
		fmt.Println(result)
	}
}
