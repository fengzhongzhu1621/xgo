package maps

import (
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/stretchr/testify/assert"
)

// 合并多个 maps 相同的 key 会被后来的 key 覆盖
func TestMerge(t *testing.T) {
	m1 := map[int]string{
		1: "a",
		2: "b",
	}
	m2 := map[int]string{
		1: "c",
		3: "d",
	}

	result := maputil.Merge(m1, m2)

	assert.Equal(t, result, map[int]string{
		1: "c",
		2: "b",
		3: "d",
	})
}
