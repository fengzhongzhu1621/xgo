package os_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunBashCommand(t *testing.T) {
	err, out, errout := RunBashCommand("echo 1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1\n", out)
	assert.Equal(t, "", errout)
}
