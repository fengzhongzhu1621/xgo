package stringutils

import (
	"testing"

	"github.com/gookit/goutil/arrutil"
	"github.com/stretchr/testify/assert"
)

func TestStringsFilter(t *testing.T) {
	is := assert.New(t)

	ss := arrutil.StringsFilter([]string{"a", "", "b", ""})
	is.Equal([]string{"a", "b"}, ss)
}
