package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

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

// 将字典转换为字符串
func TestCodec_Encode(t *testing.T) {
	codec := Codec{}

	b, err := codec.Encode(data)
	if err != nil {
		t.Fatal(err)
	}

	if encoded != string(b) {
		t.Fatalf("decoded value does not match the expected one\nactual:   %#v\nexpected: %#v", string(b), encoded)
	}
}

// 将字符串转换为字典
func TestCodec_Decode(t *testing.T) {
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

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecodeJSON(t *testing.T) {
	s := []byte(`{"name":"bob","age":18}`)
	var ts testStruct
	DecodeJSON(s, &ts)
	assert.Equal(t, "bob", ts.Name)
	assert.Equal(t, 18, ts.Age)
}

func TestDecodeJSONReader(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	// []bytes -> io.Reader
	s := []byte(`{"name":"bob","age":18}`)
	var ts testStruct
	r := bytes.NewReader(s)
	DecodeJSONReader(r, &ts)
	assert.Equal(t, "bob", ts.Name)
	assert.Equal(t, 18, ts.Age)
}

func TestEncodeJSON(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	ts := testStruct{Name: "bob", Age: 18}
	var s []byte
	EncodeJSON(ts, &s)
	assert.Equal(t, `{"name":"bob","age":18}`, string(s))
}

func TestEncJSONWriter(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	ts := testStruct{Name: "bob", Age: 18}
	var buf bytes.Buffer
	EncodeJSONWriter(ts, &buf)
	assert.Equal(t, `{"name":"bob","age":18}`, buf.String())
}

func TestInit(t *testing.T) {
	actual := reflect.TypeOf(map[string]any(nil))
	assert.Equal(t, defaultJsonHandle.MapType, actual)
}

func TestJsonMarshal(t *testing.T) {
	bk_biz_ids := []int64{1, 2, 3}
	actual, _ := json.Marshal([]map[string]any{{"field": "bk_biz_id",
		"operator": "in",
		"value":    bk_biz_ids}})
	expect := []byte(`[{"field":"bk_biz_id","operator":"in","value":[1,2,3]}]`)
	assert.Equal(t, expect, actual)
}

func TestJsonRawMessage(t *testing.T) {
	type Data struct {
		Name  string          `json:"name"`
		Value json.RawMessage `json:"value"`
	}

	jsonData := `{"name": "example", "value": {"key": "value"}}`

	var data Data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Name:", data.Name)                  // Name: example
	fmt.Printf("Value (raw JSON): %s\n", data.Value) // Value (raw JSON): {"key": "value"}

	// 如果你知道 value 的结构，可以在这里解析它
	var value map[string]string
	err = json.Unmarshal(data.Value, &value)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Value (parsed):", value) // Value (parsed): map[key:value]
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
