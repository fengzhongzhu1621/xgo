package gomonkey

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestApplyFuncReturn(t *testing.T) {
	patches := gomonkey.ApplyFuncReturn(add, 3)
	defer patches.Reset()

	assert.Equal(t, 3, add(1, 1))
}
