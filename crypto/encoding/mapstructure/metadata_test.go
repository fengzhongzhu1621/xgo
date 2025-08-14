package mapstructure

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Person4 struct {
	Name string
	Age  int
}

func TestMetadata(t *testing.T) {
	m := map[string]interface{}{
		"name": "dj",
		"age":  18,
		"job":  "programmer",
	}

	var p Person4
	var metadata mapstructure.Metadata
	mapstructure.DecodeMetadata(m, &p, &metadata)

	// keys:[]string{"Name", "Age"} unused:[]string{"job"}
	fmt.Printf("keys:%#v unused:%#v\n", metadata.Keys, metadata.Unused)
}
