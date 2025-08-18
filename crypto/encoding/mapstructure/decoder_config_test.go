package mapstructure

import (
	"fmt"
	"log"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Person7 struct {
	Name string
	Age  int
}

func TestDecoderConfig(t *testing.T) {
	m := map[string]interface{}{
		"name": 123,
		"age":  "18",
		"job":  "programmer",
	}

	var p Person7
	var metadata mapstructure.Metadata

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &p, // 输出
		Metadata:         &metadata,
	})

	if err != nil {
		log.Fatal(err)
	}

	// m -> p
	err = decoder.Decode(m)
	if err == nil {
		fmt.Println("person:", p)
		fmt.Printf("keys:%#v, unused:%#v\n", metadata.Keys, metadata.Unused)
	} else {
		fmt.Println(err.Error())
	}
}
