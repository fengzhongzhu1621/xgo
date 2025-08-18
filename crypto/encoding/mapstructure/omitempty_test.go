package mapstructure

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Person3 struct {
	Name string
	Age  int
	Job  string `mapstructure:",omitempty"`
}

func TestOmitempty(t *testing.T) {
	p := &Person3{
		Name: "dj",
		Age:  18,
	}

	// struct -> map[string]interface{}
	var m map[string]interface{}
	mapstructure.Decode(p, &m)

	// map[string]interface{} -> string
	data, _ := json.Marshal(m)
	fmt.Println(string(data)) // {"Age":18,"Name":"dj"}
}
