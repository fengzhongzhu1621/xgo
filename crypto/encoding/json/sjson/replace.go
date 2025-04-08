package sjson

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// ReplaceJsonKey replace the oldKey with newKey in jsonData
func ReplaceJsonKey(jsonData []byte, keyMap map[string]string) ([]byte, error) {
	var err error
	for oldKey, newKey := range keyMap {
		// 添加新的 key，值为旧的 key 的值
		jsonData, err = sjson.SetRawBytes(jsonData, newKey, []byte(gjson.GetBytes(jsonData, oldKey).Raw))
		if err != nil {
			return nil, fmt.Errorf("set %s key using %s value failed, err: %v", newKey, oldKey, err)
		}
		// 删除旧的key
		jsonData, err = sjson.DeleteBytes(jsonData, oldKey)
		if err != nil {
			return nil, fmt.Errorf("remove %s key failed, err: %v", oldKey, err)
		}
	}
	return jsonData, nil
}
