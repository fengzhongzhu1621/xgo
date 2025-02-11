package gjson

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

// @reverse：翻转一个数组；
// @ugly：移除 JSON 中的所有空白符；
// @pretty：使 JSON 更易用阅读；
// @this：返回当前的元素，可以用来返回根元素；
// @valid：校验 JSON 的合法性；
// @flatten：数组平坦化，即将["a", ["b", "c"]]转为["a","b","c"]；
// @join：将多个对象合并到一个对象中。
func TestModifier(t *testing.T) {
	// 读取数组children，然后使用修饰符@reverse翻转之后返回
	// ["Jack","Alex","Sara"]
	fmt.Println(gjson.Get(json2, "children|@reverse"))
	// 在上面翻转的基础上读取第一个元素，即原数组的最后一个元素
	// Jack
	fmt.Println(gjson.Get(json2, "children|@reverse|0"))
	// 移除friends数组中的所有空白字符，返回一行长长的字符串
	// [{"first":"Dale","last":"Murphy","age":44,"nets":["ig","fb","tw"]},{"first":"Roger","last":"Craig","age":68,"nets":["fb","tw"]},{"first":"Jane","last":"Murphy","age":47,"nets":["ig","tw"]}]
	fmt.Println(gjson.Get(json2, "friends|@ugly"))
	fmt.Println(gjson.Get(json2, "friends|@pretty"))
	// 返回原始的 JSON 串
	fmt.Println(gjson.Get(json2, "@this"))

	// 将数组nested的内层数组平坦到外层后返回，即将所有内层数组的元素依次添加到外层数组后面并移除内层数组
	nestedJSON := `{"nested": ["one", "two", ["three", "four"]]}`
	// ["one","two","three", "four"]
	fmt.Println(gjson.Get(nestedJSON, "nested|@flatten"))

	// 将一个数组中的各个对象合并到一个中，例子中将数组中存放的部分个人信息合并成一个对象返回
	userJSON := `{"info":[{"name":"dj", "age":18},{"phone":"123456789","email":"dj@example.com"}]}`
	// {"name":"dj","age":18,"phone":"123456789","email":"dj@example.com"}
	fmt.Println(gjson.Get(userJSON, "info|@join"))

	// 修饰符参数，通过在修饰符后加:后跟参数。
	// 可以使用@pretty修饰符的sortKeys参数对键进行排序。
	fmt.Println(gjson.Get(json2, `friends|@pretty:{"sortKeys":true}`))

	// 指定每行缩进indent（默认两个空格），每行开头字符串prefix（默认为空串）和一行最多显示字符数width（默认 80 字符）。
	// 下面在每行前增加两个空格
	fmt.Println(gjson.Get(json2, `friends|@pretty:{"sortKeys":true,"prefix":"  "}`))
}

// 自定义修饰符
func TestAddModifier(t *testing.T) {
	gjson.AddModifier("case", func(json, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(json)
		}

		if arg == "lower" {
			return strings.ToLower(json)
		}

		return json
	})

	const json = `{"children": ["Sara", "Alex", "Jack"]}`

	// ["SARA", "ALEX", "JACK"]
	fmt.Println(gjson.Get(json, "children|@case:upper"))
	// ["sara", "alex", "jack"]
	fmt.Println(gjson.Get(json, "children|@case:lower"))
}
