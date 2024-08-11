package maps

import (
	"sort"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/stretchr/testify/assert"
)

// 遍历字段，对 value 求和
func TestForEach(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	var sum int

	maputil.ForEach(m, func(_ string, value int) {
		sum += value
	})

	assert.Equal(t, sum, 10)
}

func TestKeys(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "a",
		3: "b",
		4: "c",
		5: "d",
	}

	keys := maputil.Keys(m)
	sort.Ints(keys)

	assert.Equal(t, keys, []int{1, 2, 3, 4, 5})
}
