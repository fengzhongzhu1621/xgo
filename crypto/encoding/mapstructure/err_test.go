package mapstructure

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Person5 struct {
	Name   string
	Age    int
	Emails []string
}

func TestDecodeError(t *testing.T) {
	m := map[string]interface{}{
		"name":   123,
		"age":    "bad value",
		"emails": []int{1, 2, 3},
	}

	var p Person5
	err := mapstructure.Decode(m, &p)
	if err != nil {
		// * 'Age' expected type 'int', got unconvertible type 'string', value: 'bad value'
		// * 'Emails[0]' expected type 'string', got unconvertible type 'int', value: '1'
		// * 'Emails[1]' expected type 'string', got unconvertible type 'int', value: '2'
		// * 'Emails[2]' expected type 'string', got unconvertible type 'int', value: '3'
		// * 'Name' expected type 'string', got unconvertible type 'int', value: '123'
		fmt.Println(err.Error())
	}
}
