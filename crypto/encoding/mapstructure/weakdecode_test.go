package mapstructure

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Person6 struct {
	Name   string
	Age    int
	Emails []string
}

func TestWeakDecode(t *testing.T) {
	m := map[string]interface{}{
		"name":   123,
		"age":    "18",
		"emails": []int{1, 2, 3},
	}

	var p Person6
	err := mapstructure.WeakDecode(m, &p)
	if err == nil {
		fmt.Println("person:", p)
	} else {
		fmt.Println(err.Error())
	}
}
