package gjson

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

const json3 = `
{"name": "Gilbert", "age": 61}
{"name": "Alexa", "age": 34}
{"name": "May", "age": 57}
{"name": "Deloise", "age": 44}`

// gjson提供..语法可以将多行数据看成一个数组，每行数据是一个元素
func TestLine(t *testing.T) {
	// 返回有多少行 JSON 数据
	// 4
	fmt.Println(gjson.Get(json3, "..#"))
	// 返回第一行, 0-based
	// {"name": "Alexa", "age": 34}
	fmt.Println(gjson.Get(json3, "..1"))
	// #后再接路径，表示对数组中每个元素读取后面的路径，将读取到的值组成一个新数组返回；..#.name表示读取每一行中的name字段
	// ["Gilbert","Alexa","May","Deloise"]
	fmt.Println(gjson.Get(json3, "..#.name"))
	// 括号中的内容(name="May")表示条件，所以该条含义为取name为"May"的行中的age字段
	// 57
	fmt.Println(gjson.Get(json3, `..#(name="May").age`))
}

// 遍历 JSON 行，回调返回false时遍历停止
func TestForEachLine(t *testing.T) {
	gjson.ForEachLine(json3, func(line gjson.Result) bool {
		fmt.Println("name:", gjson.Get(line.String(), "name"))
		return true
	})
}
