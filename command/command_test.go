package command

import (
	"runtime"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRunBashCommand(t *testing.T) {
	sysType := runtime.GOOS
	if sysType != "windows" {
		err, out, errout := RunBashCommand("echo 1")
		assert.Equal(t, nil, err)
		assert.Equal(t, "1\n", out)
		assert.Equal(t, "", errout)
	}
}
