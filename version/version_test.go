package version

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestLessThan(t *testing.T) {
	v1, _ := version.NewVersion("1.9")
	v2, _ := version.NewVersion("2.0")

	assert.Equal(t, true, v1.LessThan(v2))
}
