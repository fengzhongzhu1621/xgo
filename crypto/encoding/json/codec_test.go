package json

import (
	"bytes"
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
	actual := reflect.TypeOf(map[string]interface{}(nil))
	assert.Equal(t, defaultJsonHandle.MapType, actual)
}
