package sjson

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fengzhongzhu1621/xgo/tests"
)

func TestReplaceJsonKey(t *testing.T) {
	originalJSON := []byte(`{
        "name": "Alice",
        "age": 30,
        "city": "New York"
    }`)

	keyMapping := map[string]string{
		"name": "fullName",
		"city": "location",
	}

	updatedJSON, err := ReplaceJsonKey(originalJSON, keyMapping)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(updatedJSON, &result); err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	tests.PrintStruct(result)
	// 输出:
	// {
	//     "fullName": "Alice",
	//     "age": 30,
	//     "location": "New York"
	// }
}
