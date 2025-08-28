package errors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestCause(t *testing.T) {
	err := errFunc()
	if err != nil {
		err = errors.Wrap(err, "wrap error with context")

		// 获取原始错误
		originalErr := errors.Cause(err)
		fmt.Println("Original error:", originalErr) // Original error: an error occurred
	}
}
