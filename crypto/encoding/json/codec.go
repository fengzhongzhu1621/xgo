package json

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/ugorji/go/codec"
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Codec 创建了一个名为defaultJsonHandle的变量，该变量的类型是codec.JsonHandle。
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

// DecodeJSON 将字符串转换为结构体或字典
func DecodeJSON(s []byte, v any) error {
	dec := codec.NewDecoderBytes(s, &defaultJsonHandle)
	return dec.Decode(v)
}

// DecodeJSONReader 将 io.Reader 转换为结构体或字典
func DecodeJSONReader(s io.Reader, v any) error {
	dec := codec.NewDecoder(s, &defaultJsonHandle)
	return dec.Decode(v)
}

// EncodeJSON 将结构体或字典转换为字符串
func EncodeJSON(v any, s *[]byte) error {
	enc := codec.NewEncoderBytes(s, &defaultJsonHandle)
	return enc.Encode(v)
}

// EncodeJSONWriter 将结构体或字典转换为字符串并写入到 io.Writer
func EncodeJSONWriter(v any, s io.Writer) error {
	enc := codec.NewEncoder(s, &defaultJsonHandle)
	return enc.Encode(v)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// JsonStringToObject 创建了一个名为defaultJsonHandle的变量，该变量的类型是codec.JsonHandle。
// 将字符串 s 转换为 json 对象 v
// jsonStringToObject attempts to unmarshall a string as JSON into
// the object passed as pointer.
func JsonStringToObject(s string, v any) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}

// TruncateJson 将对象转换为json 字符串，并指定长度截断
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

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var iteratorJson = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	UseNumber:              true,
}.Froze()

func JsonMarshalToString(v interface{}) (string, error) {
	return iteratorJson.MarshalToString(v)
}

func JsonMarshal(v interface{}) ([]byte, error) {
	return iteratorJson.Marshal(v)
}

func JsonMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return iteratorJson.MarshalIndent(v, prefix, indent)
}

func JsonUnmarshalFromString(str string, v interface{}) error {
	return iteratorJson.UnmarshalFromString(str, v)
}

func JsonUnmarshal(data []byte, v interface{}) error {
	return iteratorJson.Unmarshal(data, v)
}

func JsonUnmarshalArray(items []string, result interface{}) error {
	strArrJSON := "[" + strings.Join(items, ",") + "]"
	return iteratorJson.Unmarshal([]byte(strArrJSON), result)
}

// ExtractFieldsFromJSON 从原始 JSON 中提取指定字段并返回新的 JSON 字符串
func ExtractFieldsFromJSON(jsonData *string, fields []string) (*string, error) {
	if jsonData == nil {
		return nil, fmt.Errorf("jsonData is nil")
	}
	if len(fields) == 0 || *jsonData == "" {
		return jsonData, nil
	}

	elements := gjson.GetMany(*jsonData, fields...)
	if len(elements) != len(fields) {
		return nil, fmt.Errorf("mismatch in number of fields and extracted elements")
	}

	result := make(map[string]interface{}, len(fields))
	for idx, field := range fields {
		if elements[idx].Exists() {
			// 根据需要处理不同的数据类型
			result[field] = elements[idx].Value()
		} else {
			result[field] = nil
		}
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return stringPtr(string(jsonBytes)), nil
}

// stringPtr 辅助函数，将字符串转换为 *string
func stringPtr(s string) *string {
	return &s
}
