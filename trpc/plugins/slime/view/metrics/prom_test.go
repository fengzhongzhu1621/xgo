package metrics

import (
	"testing"
)

// Prometheus cannot properly test in UT.
// We just try to call exported functions and make sure they are not panic.
// Note, the number of tag pairs must match metrics.TagNamesXxx.
func TestEmitter(t *testing.T) {
	tagsApp := []string{
		TagCaller, "caller_",
		TagCallee, "callee_",
		TagMethod, "method_",
		TagAttempts, "2",
		TagErrCodes, "0",
		TagThrottled, "false",
		TagInflight, "1",
		TagNoMoreAttempt, "false",
	}
	tagsReal := []string{
		TagCaller, "caller_",
		TagCallee, "callee_",
		TagMethod, "method_",
		TagErrCodes, "123",
		TagInflight, "false",
		TagNoMoreAttempt, "true",
	}
	m := NewEmitter()
	m.Inc(FQNAppRequest, 1, tagsApp...)
	m.Inc(FQNRealRequest, 1, tagsReal...)
	m.Observe(FQNAppCostMs, 10, tagsApp...)
	m.Observe(FQNRealCostMs, 10, tagsReal...)
}
