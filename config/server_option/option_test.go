package server_option

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestNewServerOption(t *testing.T) {
	op := NewServerOption()
	op.AddFlags(pflag.CommandLine)
}
