package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSONMerge(t *testing.T) {
	dataJSON := `{
        "name": "Alice",
        "age": 30,
        "address": {
            "city": "Wonderland"
        }
    }`

	// 原有的数据保持不变
	// 同属性覆盖
	// 新增的属性添加
	patchJSON := `{
        "age": 31,
        "email": "alice@example.com"
    }`

	var data, patch json.RawMessage
	json.Unmarshal([]byte(dataJSON), &data)
	json.Unmarshal([]byte(patchJSON), &patch)

	merged, err := JSONMerge(data, patch)
	if err != nil {
		fmt.Println("Error merging JSON:", err)
		return
	}

	fmt.Println(
		string(merged),
	) // {"address":{"city":"Wonderland"},"age":31,"email":"alice@example.com","name":"Alice"}
}
