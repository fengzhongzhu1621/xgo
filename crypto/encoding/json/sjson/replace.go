package sjson

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// ReplaceJsonKey 遍历 keyMap 中的每一个键值对，将 JSON 数据中的旧键 (oldKey) 替换为新键 (newKey)
func ReplaceJsonKey(jsonData []byte, keyMap map[string]string) ([]byte, error) {
	var err error

	// 一次性解析 JSON
	result := gjson.ParseBytes(jsonData)

	for oldKey, newKey := range keyMap {
		// 检查旧键是否存在
		value := gjson.GetBytes(jsonData, oldKey)
		if !value.Exists() {
			return nil, fmt.Errorf("key %s does not exist in JSON", oldKey)
		}

		// 检查新键是否已存在（根据需求决定是否覆盖）
		if result.Get(newKey).Exists() {
			return nil, fmt.Errorf("key %s already exists in JSON", newKey)
		}

		// 删除旧的 key
		jsonData, err = sjson.DeleteBytes(jsonData, oldKey)
		if err != nil {
			return nil, fmt.Errorf("remove %s key failed, err: %v", oldKey, err)
		}

		// 添加新的 key，值为旧的 key 的值
		jsonData, err = sjson.SetRawBytes(jsonData, newKey, []byte(value.Raw))
		if err != nil {
			return nil, fmt.Errorf("set %s key using %s value failed, err: %v", newKey, oldKey, err)
		}
	}

	return jsonData, nil
}
