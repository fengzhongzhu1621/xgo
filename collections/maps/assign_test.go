package maps

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Assign(map[string]int{"a": 1, "b": 2}, map[string]int{"b": 3, "c": 4})

	is.Len(result1, 3)
	is.Equal(result1, map[string]int{"a": 1, "b": 3, "c": 4})

	type myMap map[string]int
	before := myMap{"": 0, "foobar": 6, "baz": 3}
	after := lo.Assign(before, before)
	is.IsType(after, before, "type preserved")
}
