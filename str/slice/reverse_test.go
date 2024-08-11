package slice

import (
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	slice1 := []int{1, 2, 3}

	reversed := utils.Reverse(slice1)
	assert.Equal(t, reversed, []int{3, 2, 1})
}
