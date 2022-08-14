package log

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

type fieldKey string

// FieldMap allows customization of the key names for default fields.
// 字段映射字典，字典的默认值为指定的key
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	// key存在返回value
	if k, ok := f[key]; ok {
		return k
	}
	// key不存在返回key，即将key作为value返回
	return string(key)
}

