package mapstructure

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Person2 struct {
	Name  string
	Age   int
	Job   string
	Other map[string]interface{} `mapstructure:",remain"`
}

func TestRemain(t *testing.T) {
	data := `
    {
      "name": "dj",
      "age":18,
      "job":"programmer",
      "height":"1.8m",
      "handsome": true
    }
  `

	var m map[string]interface{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatal(err)
	}

	var p Person2
	mapstructure.Decode(m, &p)
	fmt.Println("other", p.Other) // map[handsome:true height:1.8m]
}
