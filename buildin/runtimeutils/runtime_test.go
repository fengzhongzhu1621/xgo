package runtimeutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPackage(t *testing.T) {
	actual := GetPackage()
	expect := "xgo/utils/runtimeutils"
	assert.Equal(t, expect, actual)
}
