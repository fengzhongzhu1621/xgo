package stringutils

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestRuneLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.Equal(5, lo.RuneLength("hellô"))
	is.Equal(6, len("hellô"))
}
