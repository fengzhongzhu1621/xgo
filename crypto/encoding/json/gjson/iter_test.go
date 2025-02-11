package gjson

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

const json4 = `
{
  "name":"dj",
  "age":18,
  "pets": ["cat", "dog"],
  "contact": {
    "phone": "123456789",
    "email": "dj@example.com"
  }
}`

// gjson还提供了通用的遍历数组和对象的方式。
// gjson.Get()方法返回一个gjson.Result类型的对象，json.Result提供了ForEach()方法用于遍历。
// 该方法接受一个类型为func (key, value gjson.Result) bool的回调函数。
// 遍历对象时key和value分别为对象的键和值；
// 遍历数组时，value为数组元素，key为空（不是索引）。回调返回false时，遍历停止。
func TestForEach(t *testing.T) {
	// 遍历数组
	pets := gjson.Get(json4, "pets")
	pets.ForEach(func(_, pet gjson.Result) bool {
		fmt.Println(pet)
		return true
	})

	// 遍历对象
	contact := gjson.Get(json4, "contact")
	contact.ForEach(func(key, value gjson.Result) bool {
		fmt.Println(key, value)
		return true
	})
}
