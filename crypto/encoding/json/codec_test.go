package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/stretchr/testify/assert"
)

// encoded form of the data
const encoded = `{
  "key": "value",
  "list": [
    "item1",
    "item2",
    "item3"
  ],
  "map": {
    "key": "value"
  },
  "nested_map": {
    "map": {
      "key": "value",
      "list": [
        "item1",
        "item2",
        "item3"
      ]
    }
  }
}`

// Viper's internal representation
var data = map[string]interface{}{
	"key": "value",
	"list": []interface{}{
		"item1",
		"item2",
		"item3",
	},
	"map": map[string]interface{}{
		"key": "value",
	},
	"nested_map": map[string]interface{}{
		"map": map[string]interface{}{
			"key": "value",
			"list": []interface{}{
				"item1",
				"item2",
				"item3",
			},
		},
	},
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	// 如果你想要在序列化时忽略某个字段，但在反序列化时仍然能够处理它，你可以使用 json.RawMessage 类型。
	Password json.RawMessage `json:"password,omitempty"` // 如果 Password 为空，则不会序列化
}

type Person2 struct {
	Name string `json:"name"`
	// 如果你将结构体中的字段定义为指针类型，并且在序列化时该字段为 nil，那么在 JSON 序列化时，该字段也会被忽略（前提是你没有使用 omitempty 以外的选项）。
	Age      *int    `json:"age,omitempty"`      // 如果 Age 为 nil，则不会序列化
	Password *string `json:"password,omitempty"` // 如果 Password 为 nil，则不会序列化
}

func TestInit(t *testing.T) {
	actual := reflect.TypeOf(map[string]any(nil))
	assert.Equal(t, defaultJsonHandle.MapType, actual)
}

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 将字典转换为字符串
func TestMapToJson(t *testing.T) {
	{
		codec := Codec{}

		b, err := codec.Encode(data)
		if err != nil {
			t.Fatal(err)
		}

		if encoded != string(b) {
			t.Fatalf("decoded value does not match the expected one\nactual:   %#v\nexpected: %#v", string(b), encoded)
		}
	}

	{
		bizIds := []int64{1, 2, 3}
		actual, _ := json.Marshal([]map[string]any{{"field": "bk_biz_id",
			"operator": "in",
			"value":    bizIds}})
		expect := []byte(`[{"field":"bk_biz_id","operator":"in","value":[1,2,3]}]`)
		assert.Equal(t, expect, actual)
	}

	{
		// TestToJson 字典转换为字符串
		// func ToJson(value any) (string, error)
		aMap := map[string]int{"a": 1, "b": 2, "c": 3}
		result, err := convertor.ToJson(aMap)

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Println(result)

		// Output:
		// {"a":1,"b":2,"c":3}
	}
}

// 将字符串转换为字典
func TestJsonToMap(t *testing.T) {
	{
		t.Run("OK", func(t *testing.T) {
			codec := Codec{}

			v := map[string]interface{}{}

			err := codec.Decode([]byte(encoded), v)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(data, v) {
				t.Fatalf("decoded value does not match the expected one\nactual:   %#v\nexpected: %#v", v, data)
			}
		})

		t.Run("InvalidData", func(t *testing.T) {
			codec := Codec{}

			v := map[string]interface{}{}

			err := codec.Decode([]byte(`invalid data`), v)
			if err == nil {
				t.Fatal("expected decoding to fail")
			}

			t.Logf("decoding failed as expected: %s", err)
		})
	}
}

func TestJsonToStruct(t *testing.T) {
	{
		jsonStr := `{"name": "Bob", "age": 25}`

		var user testStruct
		err := json.Unmarshal([]byte(jsonStr), &user)
		if err != nil {
			fmt.Println("转换错误:", err)
			return
		}

		fmt.Println(user) // 输出：{Bob 25}
	}

	{
		s := []byte(`{"name":"bob","age":18}`)
		var ts testStruct
		DecodeJSON(s, &ts)
		assert.Equal(t, "bob", ts.Name)
		assert.Equal(t, 18, ts.Age)
	}

	{
		// []bytes -> io.Reader
		s := []byte(`{"name":"bob","age":18}`)
		var ts testStruct
		r := bytes.NewReader(s)
		DecodeJSONReader(r, &ts)
		assert.Equal(t, "bob", ts.Name)
		assert.Equal(t, 18, ts.Age)
	}
}

func TestStructToJson(t *testing.T) {
	{
		ts := testStruct{Name: "bob", Age: 18}
		var s []byte
		EncodeJSON(ts, &s)
		assert.Equal(t, `{"name":"bob","age":18}`, string(s))
	}

	{
		ts := testStruct{Name: "bob", Age: 18}
		var buf bytes.Buffer
		EncodeJSONWriter(ts, &buf)
		assert.Equal(t, `{"name":"bob","age":18}`, buf.String())
	}
}

func TestOmitEmpty(t *testing.T) {
	p := Person{
		Name: "Alice",
		Age:  30,
		// Password 字段为空，因此在序列化时会被忽略
	}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes)) // 输出: {"name":"Alice","age":30}

	p2 := Person2{
		Name: "Alice",
		// Age 和 Password 字段为 nil，因此在序列化时会被忽略
	}

	jsonBytes2, err := json.Marshal(p2)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes2)) // 输出: {"name":"Alice"}
}
