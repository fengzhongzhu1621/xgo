package maps

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestInvert(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.Invert(map[string]int{"a": 1, "b": 2})
	r2 := lo.Invert(map[string]int{"a": 1, "b": 2, "c": 1})

	is.Len(r1, 2)
	is.EqualValues(map[int]string{1: "a", 2: "b"}, r1)
	is.Len(r2, 2)
}
