package gjson

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const json = `{
  "code": 200,
  "data":{
    "user": {
      "id": 1,
      "name": "name_1",
    }
  }
}`

const json2 = `
{
  "name":{"first":"Tom", "last": "Anderson"},
  "age": 37,
  "children": ["Sara", "Alex", "Jack"],
  "fav.movie": "Dear Hunter",
  "friends": [
    {"first": "Dale", "last":"Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`

// gjson.Get()函数实际上返回的是gjson.Result类型，需要调用其相应的方法进行转换对应的类型
func TestGet(t *testing.T) {
	code := gjson.Get(json, "code").Int()
	assert.Equal(t, code, int64(200))

	userID := gjson.Get(json, "data.user.id").String()
	assert.Equal(t, userID, "1")

	userID2 := gjson.Get(json, "data.user.id").Int()
	assert.Equal(t, userID2, int64(1))
}

// gjson支持在键中包含通配符*和?，*匹配任意多个字符，?匹配单个字符，例如ca*可以匹配cat/cate/cake等以ca开头的键，ca?只能匹配cat/cap等以ca开头且后面只有一个字符的键。
// 数组使用键名 + . + 索引（索引从 0 开始）的方式读取元素，如果键pets对应的值是一个数组，那么pets.0读取数组的第一个元素，pets.1读取第二个元素。
// 数组长度使用**键名 + . + #**获取，例如pets.#返回数组pets的长度。
// 如果键名中出现.，那么需要使用\进行转义。
//
// children.#：返回数组children的长度；
// children.1：读取数组children的第 2 个元素（注意索引从 0 开始）；
// child*.2：首先child*匹配children，.2读取第 3 个元素；
// c?ildren.0：c?ildren匹配到children，.0读取第一个元素；
// fav\.movie：因为键名中含有.，故需要\转义；
// friends.#.first：如果数组后#后还有内容，则以后面的路径读取数组中的每个元素，返回一个新的数组。所以该查询返回的数组所有friends的first字段组成；
// friends.1.last：读取friends第 2 个元素的last字段。
func TestKeyPath(t *testing.T) {
	// last name: Anderson
	fmt.Println("last name:", gjson.Get(json2, "name.last"))
	// age: 37
	fmt.Println("age:", gjson.Get(json2, "age"))
	// children: ["Sara", "Alex", "Jack"]
	fmt.Println("children:", gjson.Get(json2, "children"))
	// children count: 3
	fmt.Println("children count:", gjson.Get(json2, "children.#"))
	// second child: Alex
	fmt.Println("second child:", gjson.Get(json2, "children.1"))
	// third child*: Jack
	fmt.Println("third child*:", gjson.Get(json2, "child*.2"))
	// first c?ild: Sara
	fmt.Println("first c?ild:", gjson.Get(json2, "c?ildren.0"))
	// fav.moive Dear Hunter
	fmt.Println("fav.moive", gjson.Get(json2, `fav\.movie`))
	// first name of friends: ["Dale","Roger","Jane"]
	fmt.Println("first name of friends:", gjson.Get(json2, "friends.#.first"))
	// last name of second friend: Craig
	fmt.Println("last name of second friend:", gjson.Get(json2, "friends.1.last"))
}

// 支持按条件查询元素，#(条件)返回第一个满足条件的元素，#(条件)#返回所有满足条件的元素。
// 括号内的条件可以有==、!=、<、<=、>、>=，还有简单的模式匹配%（符合某个模式），!%（不符合某个模式）
func TestCondition(t *testing.T) {
	// friends.#(last="Murphy")返回数组friends中第一个last为Murphy的元素，.first表示取出该元素的first字段返回
	// Dale
	fmt.Println(gjson.Get(json2, `friends.#(last="Murphy").first`))
	// friends.#(last="Murphy")#返回数组friends中所有的last为Murphy的元素，然后读取它们的first字段放在一个数组中返回。
	// ["Dale","Jane"]
	fmt.Println(gjson.Get(json2, `friends.#(last="Murphy")#.first`))
	// friends.#(age>45)#返回数组friends中所有年龄大于 45 的元素，然后读取它们的last字段返回
	// ["Craig","Murphy"]
	fmt.Println(gjson.Get(json2, "friends.#(age>45)#.last"))
	// friends.#(first%"D*")返回数组friends中第一个first字段满足模式D*的元素，取出其last字段返回
	// Murphy
	fmt.Println(gjson.Get(json2, `friends.#(first%"D*").last`))
	// 返回数组friends中第一个first字段**不**满足模式D*的元素，读取其last`字段返回
	// Craig
	fmt.Println(gjson.Get(json2, `friends.#(first!%"D*").last`))
	// 这是个嵌套条件，friends.#(nets.#(=="fb"))#返回数组friends的元素的nets字段中有fb的所有元素，然后取出first字段返回。
	// ["Dale","Roger"]
	fmt.Println(gjson.Get(json2, `friends.#(nets.#(=="fb"))#.first`))
}
