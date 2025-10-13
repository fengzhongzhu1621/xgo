package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	require := require.New(t)
	p := &FileProvider{}
	require.Equal("file", p.Name())
	RegisterProvider(p)
	pp := GetProvider("file")
	require.Equal(p, pp)
}
