package xerror

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNewErrFrequency(t *testing.T) {
	ef := NewErrFrequency(errors.New("example error"), 1*time.Second)

	fmt.Println(ef.IsErrAlwaysAppear(errors.New("example error"))) // false
	time.Sleep(1 * time.Second)

	fmt.Println(ef.IsErrAlwaysAppear(errors.New("example error"))) // true
	ef.Release()

	fmt.Println(ef.IsErrAlwaysAppear(errors.New("example error"))) // false
}
