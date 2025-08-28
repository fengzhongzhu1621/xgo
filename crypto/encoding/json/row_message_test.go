package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonRawMessage(t *testing.T) {
	type Data struct {
		Name  string          `json:"name"`
		Value json.RawMessage `json:"value"`
	}

	jsonData := `{"name": "example", "value": {"key": "value"}}`

	var data Data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Name:", data.Name)                  // Name: example
	fmt.Printf("Value (raw JSON): %s\n", data.Value) // Value (raw JSON): {"key": "value"}

	// 如果你知道 value 的结构，可以在这里解析它
	var value map[string]string
	err = json.Unmarshal(data.Value, &value)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Value (parsed):", value) // Value (parsed): map[key:value]
}

type RawMessageSlice struct {
	Name       string                       `json:"name"`
	UpsertInfo map[string][]json.RawMessage `json:"upsert_info"`
}

func TestJsonRawMessageSlice(t *testing.T) {
	// 示例 JSON 数据
	jsonData := `{
        "name": "transaction_1",
        "upsert_info": {
            "users": [
                {"id": 1, "name": "Alice"},
                {"id": 2, "name": "Bob"}
            ],
            "products": [
                {"id": 101, "name": "Laptop"},
                {"id": 102, "name": "Phone"}
            ]
        }
    }`

	var data RawMessageSlice
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	fmt.Printf("Parsed Data: %+v\n", data)

	// 访问特定数据
	for key, messages := range data.UpsertInfo {
		fmt.Printf("Upserting %s:\n", key)
		for _, msg := range messages {
			fmt.Println(string(msg))
		}
	}

	// Outputs:
	//
	// 	Upserting users:
	// {"id": 1, "name": "Alice"}
	// {"id": 2, "name": "Bob"}
	// Upserting products:
	// {"id": 101, "name": "Laptop"}
	// {"id": 102, "name": "Phone"}

	// 编码回 JSON
	encodedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println("Encoded JSON:")
	fmt.Println(string(encodedJSON))
	// Outputs
	//
	// {
	// 	"name": "transaction_1",
	// 	"upsert_info": {
	// 	  "products": [
	// 		{
	// 		  "id": 101,
	// 		  "name": "Laptop"
	// 		},
	// 		{
	// 		  "id": 102,
	// 		  "name": "Phone"
	// 		}
	// 	  ],
	// 	  "users": [
	// 		{
	// 		  "id": 1,
	// 		  "name": "Alice"
	// 		},
	// 		{
	// 		  "id": 2,
	// 		  "name": "Bob"
	// 		}
	// 	  ]
	// 	}
	//   }
}
