package time_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRetryBackoff(t *testing.T) {
	for i := 0; i <= 16; i++ {
		backoff := RetryBackoff(i, time.Millisecond, 512*time.Millisecond)
		assert.Equal(t, true, backoff >= 0)
		assert.Equal(t, true, backoff <= 512*time.Millisecond)
	}
}