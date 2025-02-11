package gjson

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

func TestValid(t *testing.T) {
	const json = `{"name":dj,age:18}`
	// gjson假设我们传入的 JSON 串是合法的。如果 JSON 非法也不会panic，这时会返回不确定的结果
	fmt.Println(gjson.Get(json, "name"))

	if !gjson.Valid(json) {
		// error
		fmt.Println("error")
	} else {
		fmt.Println("ok")
	}
}
