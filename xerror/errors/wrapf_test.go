package errors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestWrapf(t *testing.T) {
	err := errFunc()
	if err != nil {
		err = errors.Wrapf(err, "error in someFunc: %s", "additional context")
		fmt.Println(err) // error in someFunc: additional context: an error occurred
	}
}
