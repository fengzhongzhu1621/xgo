package json

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/ugorji/go/codec"
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 创建了一个名为defaultJsonHandle的变量，该变量的类型是codec.JsonHandle。
type Codec struct{}

// Encode 将字典转换为字符串
func (Codec) Encode(v map[string]interface{}) ([]byte, error) {
	// TODO: expose prefix and indent in the Codec as setting?
	return json.MarshalIndent(v, "", "  ")
}

// Decode 将字符串转换为字典
func (Codec) Decode(b []byte, v map[string]interface{}) error {
	return json.Unmarshal(b, &v)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 创建了一个名为defaultJsonHandle的变量，该变量的类型是codec.JsonHandle。
// codec.JsonHandle是Go标准库encoding/json包中的一个结构体，用于自定义JSON编码和解码的行为。
//
// MapKeyAsString字段被设置为true。这意味着在解码JSON时，所有的映射（map）键将被解析为字符串类型，而不是默认的interface{}类型。
// 这可以避免在处理JSON数据时出现类型断言的问题。
var defaultJsonHandle = codec.JsonHandle{MapKeyAsString: true}

func init() {
	// 创建了一个类型为 map[string]interface{} 的空映射（map），并将其初始化为 nil。
	defaultJsonHandle.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

// DecJSON 将字符串转换为结构体或字典
func DecodeJSON(s []byte, v any) error {
	dec := codec.NewDecoderBytes(s, &defaultJsonHandle)
	return dec.Decode(v)
}

// DecJSONReader 将 io.Reader 转换为结构体或字典
func DecodeJSONReader(s io.Reader, v any) error {
	dec := codec.NewDecoder(s, &defaultJsonHandle)
	return dec.Decode(v)
}

// EncJSON 将结构体或字典转换为字符串
func EncodeJSON(v any, s *[]byte) error {
	enc := codec.NewEncoderBytes(s, &defaultJsonHandle)
	return enc.Encode(v)
}

// EncJSONWriter 将结构体或字典转换为字符串并写入到 io.Writer
func EncodeJSONWriter(v any, s io.Writer) error {
	enc := codec.NewEncoder(s, &defaultJsonHandle)
	return enc.Encode(v)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 创建了一个名为defaultJsonHandle的变量，该变量的类型是codec.JsonHandle。
// 将字符串 s 转换为 json 对象 v
// jsonStringToObject attempts to unmarshall a string as JSON into
// the object passed as pointer.
func JsonStringToObject(s string, v any) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}

// Truncate 将对象转换为json 字符串，并指定长度截断
func TruncateJson(args any, length int) string {
	s, err := jsoniter.MarshalToString(args)
	if err != nil {
		s = fmt.Sprintf("%v", args)
	}

	if length > len(s) {
		return s
	}
	return s[:length]
}
