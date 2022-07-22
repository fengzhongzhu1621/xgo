package log

import (
	"testing"
)

func TestLogrusInfo(t *testing.T) {
	LogrusSetLevel(TraceLevel)
	LogrusInfo("info msg")
}
