package stringutils

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/formatter"
)

// Pretty data to JSON string.
// func Pretty(v any) (string, error)

func TestPretty(t *testing.T) {
	result1, _ := formatter.Pretty([]string{"a", "b", "c"})
	result2, _ := formatter.Pretty(map[string]int{"a": 1})

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// [
	//     "a",
	//     "b",
	//     "c"
	// ]
	// {
	//     "a": 1
	// }
}

// Pretty encode data to writer.
// func PrettyToWriter(v any, out io.Writer) error
func TestPrettyToWriter(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Aage uint   `json:"age"`
	}
	user := User{Name: "King", Aage: 10000}

	buf := &bytes.Buffer{}
	err := formatter.PrettyToWriter(user, buf)

	fmt.Println(buf)
	fmt.Println(err)

	// Output:
	// {
	//     "name": "King",
	//     "age": 10000
	// }
	//
	// <nil>
}
