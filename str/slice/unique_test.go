package slice

import (
	"fmt"
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
)

// 去重切片
func TestUnique(t *testing.T) {
	withDuplicates := []int{1, 2, 2, 3, 3, 3, 4}
	unique := utils.Unique(withDuplicates)
	fmt.Println(unique) // Output: [1 2 3 4]
}
