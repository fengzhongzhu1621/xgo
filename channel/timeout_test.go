package channel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeoutCaller(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "bad function",
			f: func(t *testing.T) {
				timeout := time.Second
				actualTime := timeout * 2
				f := func(errChan chan error) {
					defer close(errChan)
					time.Sleep(actualTime)
					errChan <- nil
				}
				now := time.Now()
				err := TimeoutCaller(f, timeout)
				assert.Error(t, err)
				// 断言两个时间点之间的差异是否在指定的时间范围内。
				assert.WithinDuration(t, now.Add(timeout), time.Now(), 250*time.Millisecond)
			},
		},
		{
			name: "common function",
			f: func(t *testing.T) {
				timeout := time.Second
				actualTime := timeout / 2
				f := func(errChan chan error) {
					defer close(errChan)
					time.Sleep(actualTime)
					errChan <- nil
				}
				now := time.Now()
				err := TimeoutCaller(f, timeout)
				assert.NoError(t, err)
				assert.WithinDuration(t, now.Add(actualTime), time.Now(), 250*time.Millisecond)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f(t)
		})
	}
}
