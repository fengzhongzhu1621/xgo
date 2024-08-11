package slice

import (
	"fmt"
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
)

func TestConcat(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}

	concatenated := utils.Concat(slice1, slice2)
	fmt.Println(concatenated) // Output: [1 2 3 4 5 6]
}
