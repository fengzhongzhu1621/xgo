package system

import (
	"testing"

	"github.com/duke-git/lancet/v2/system"
	"github.com/stretchr/testify/assert"
)

// func CompareOsEnv(key, comparedEnv string) bool
func TestCompareOsEnv(t *testing.T) {
	err := system.SetOsEnv("foo", "abc")
	if err != nil {
		return
	}

	result := system.CompareOsEnv("foo", "abc")
	assert.Equal(t, true, result)
}
