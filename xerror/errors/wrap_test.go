package errors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func errFunc() error {
	return errors.New("an error occurred")
}

// 将错误进行包装，并附加额外的上下文信息。
func TestWrap(t *testing.T) {
	err := errFunc()
	if err != nil {
		err = errors.Wrap(err, "while calling someFunc")
		fmt.Println(err) // while calling someFunc: an error occurred
	}
}
