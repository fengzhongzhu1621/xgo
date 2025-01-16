package stringutils

import (
	"testing"

	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"
)

func TestChunkString(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.ChunkString("12345", 2)
	is.Equal([]string{"12", "34", "5"}, result1)

	result2 := lo.ChunkString("123456", 2)
	is.Equal([]string{"12", "34", "56"}, result2)

	result3 := lo.ChunkString("123456", 6)
	is.Equal([]string{"123456"}, result3)

	result4 := lo.ChunkString("123456", 10)
	is.Equal([]string{"123456"}, result4)

	result5 := lo.ChunkString("", 2)
	is.Equal([]string{""}, result5)

	result6 := lo.ChunkString("明1好休2林森", 2)
	is.Equal([]string{"明1", "好休", "2林", "森"}, result6)

	is.Panics(func() {
		lo.ChunkString("12345", 0)
	})
}
