package xgo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestVersion(t *testing.T) {
	actual := Version()
	expect := fmt.Sprintf("%s (%s %s-%s) %s",
		version, runtime.Version(), runtime.GOOS, runtime.GOARCH, buildInfo)
	// 0.0.1 (go1.15.5 darwin-amd64) -----
	// fmt.Print(expect)
	assert.Equal(t, expect, actual)
}
