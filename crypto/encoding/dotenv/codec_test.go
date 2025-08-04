package dotenv

import (
	"reflect"
	"testing"
)

// original form of the data
const original = `# key-value pair
KEY=value
`

// encoded form of the data
const encoded = `KEY=value
`

// Viper's internal representation
var data = map[string]interface{}{
	"KEY": "value",
}

func TestCodec_Encode(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		codec := Codec{}

		b, err := codec.Encode(data)
		if err != nil {
			t.Fatal(err)
		}

		if encoded != string(b) {
			t.Fatalf(
				"decoded value does not match the expected one\nactual:   %#v\nexpected: %#v",
				string(b),
				encoded,
			)
		}
	})

	t.Run("flatten map", func(t *testing.T) {
		codec := Codec{}

		data2 := map[string]interface{}{
			"KEY": map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		}
		b, err := codec.Encode(data2)
		if err != nil {
			t.Fatal(err)
		}

		actual := `KEY_A=1
KEY_B=2
`
		if string(b) != actual {
			t.Fatalf("decoded value does not match the expected one\nactual:   %#v\nexpected: %#v",
				string(b), actual)
		}
	})
}

func TestCodec_Decode(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		codec := Codec{}

		v := map[string]interface{}{}

		err := codec.Decode([]byte(original), v)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(data, v) {
			t.Fatalf(
				"decoded value does not match the expected one\nactual:   %#v\nexpected: %#v",
				v,
				data,
			)
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
