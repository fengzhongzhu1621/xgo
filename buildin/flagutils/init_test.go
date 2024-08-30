package flagutils

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestInitFlag(t *testing.T) {
	InitFlags()
	pflag.PrintDefaults()
}
