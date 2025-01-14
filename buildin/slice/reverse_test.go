package slice

import (
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	slice1 := []int{1, 2, 3}

	reversed := utils.Reverse(slice1)
	assert.Equal(t, []int{3, 2, 1}, reversed)
}

func TestLancetReverse(t *testing.T) {
	strs := []string{"a", "b", "c", "d"}

	slice.Reverse(strs)
	assert.Equal(t, []string{"d", "c", "b", "a"}, strs)
}
