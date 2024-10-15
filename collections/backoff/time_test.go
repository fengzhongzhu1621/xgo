package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryBackoff(t *testing.T) {
	for i := 0; i <= 16; i++ {
		backoff, _ := RetryBackoff(i, time.Millisecond, 512*time.Millisecond)
		assert.Equal(t, true, backoff >= 0)
		assert.Equal(t, true, backoff <= 512*time.Millisecond)
	}
}
