package cast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr2map(t *testing.T) {
	s := "a=1&b=2&c="
	actual := Str2map(s, "&", "=")
	expect := map[string]string{"a": "1", "b": "2", "c": ""}
	assert.Equal(t, expect, actual)
}
