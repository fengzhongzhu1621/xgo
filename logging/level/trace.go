package level

import (
	"os"

	"github.com/fengzhongzhu1621/xgo/env"
)

var traceEnabled = traceEnableFromEnv()

// traceEnableFromEnv checks whether trace is enabled by reading from environment.
// Close trace if empty or zero, open trace if not zero, default as closed.
func traceEnableFromEnv() bool {
	switch os.Getenv(env.LogTrace) {
	case "":
		fallthrough
	case "0":
		return false
	default:
		return true
	}
}

// EnableTrace enables trace.
func EnableTrace() {
	traceEnabled = true
}
