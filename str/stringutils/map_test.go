package stringutils

import (
	"testing"

	"github.com/gookit/goutil/arrutil"
	"github.com/stretchr/testify/assert"
)

func TestStringsMap(t *testing.T) {
	is := assert.New(t)

	ss := arrutil.StringsMap([]string{"a", "b", "c"}, func(s string) string {
		return s + "1"
	})
	is.Equal([]string{"a1", "b1", "c1"}, ss)
}
