package tests

import (
	"fmt"

	"github.com/dablelv/cyan/encoding"
)

func PrintStruct(obj interface{}) {
	s, _ := encoding.ToIndentJSON(obj)
	fmt.Printf("%v\n", s)
}
