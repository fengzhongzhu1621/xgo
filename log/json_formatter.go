package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
)

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    },
	// }
	// 字段映射表，将
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	// 将entry中的error转换为字符串，放到字典data中
	// 在初始化时，可以提前定义好 map 所需要的容量（空间大小），当添加的键值超过容量时自动加一。
	data := make(Fields, len(entry.Data)+4) // 指定字典的容量
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	// 将data放到一个指定key的字典中
	if f.DataKey != "" {
		newData := make(Fields, 4)
		newData[f.DataKey] = data
		data = newData
	}

	// 给指定字段加上前缀，防止和预留字段冲突
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())

	// 选择时间格式
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	// 添加错误信息
	if entry.err != "" {
		data[f.FieldMap.resolve(FieldKeyLogrusError)] = entry.err
	}

	// 设置日志时间，时间默认是日志的创建时间
	if !f.DisableTimestamp {
		data[f.FieldMap.resolve(FieldKeyTime)] = entry.Time.Format(timestampFormat)
	}
	// 设置日志内容
	data[f.FieldMap.resolve(FieldKeyMsg)] = entry.Message
	// 设置日志级别
	data[f.FieldMap.resolve(FieldKeyLevel)] = entry.Level.String()

	// 设置调用堆栈信息，例如函数名和行数
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
		if funcVal != "" {
			data[f.FieldMap.resolve(FieldKeyFunc)] = funcVal
		}
		if fileVal != "" {
			data[f.FieldMap.resolve(FieldKeyFile)] = fileVal
		}
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 将字典转换为json字符串
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!f.DisableHTMLEscape)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return b.Bytes(), nil
}
