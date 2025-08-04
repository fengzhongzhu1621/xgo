package tuple

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUnzip(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1, r2 := lo.Unzip2([]lo.Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}})

	is.Equal(r1, []string{"a", "b"})
	is.Equal(r2, []int{1, 2})
}

func TestUnzipBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1, r2 := lo.UnzipBy2(
		[]lo.Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}},
		func(i lo.Tuple2[string, int]) (a string, b int) {
			return i.A + i.A, i.B + i.B
		},
	)

	is.Equal(r1, []string{"aa", "bb"})
	is.Equal(r2, []int{2, 4})
}
