package newset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestAppend(*testing.T) {
	s := mapset.NewSet[int]()
	s.Append([]int{1, 2, 3}...)
	fmt.Println(s) // Set{2, 3, 1}
}
