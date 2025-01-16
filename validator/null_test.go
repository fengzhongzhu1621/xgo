package validator

import (
	"testing"

	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	var x int
	is.False(lo.IsNil(x))

	var k struct{}
	is.False(lo.IsNil(k))

	var s *string
	is.True(lo.IsNil(s))

	var i *int
	is.True(lo.IsNil(i))

	var b *bool
	is.True(lo.IsNil(b))

	var ifaceWithNilValue any = (*string)(nil) //nolint:staticcheck
	is.True(lo.IsNil(ifaceWithNilValue))
	is.False(ifaceWithNilValue == nil) //nolint:staticcheck
}
