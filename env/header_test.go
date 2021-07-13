package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppendEnv(t *testing.T) {
	env := []string{"a", "b"}
	key := "key"
	value1 := "value1"
	value2 := "value2"
	env1 := AppendEnv(env, key, value1, value2)
	assert.Equal(t, []string{"a", "b"}, env)
	assert.Equal(t, []string{"a", "b", "KEY=value1, value2"}, env1)
}
