package errors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestPrintStack(t *testing.T) {
	err := errFunc()
	if err != nil {
		err = errors.Wrap(err, "something went wrong")
		// Print the error with the stack trace
		fmt.Printf("Error: %+v\n", err)
	}
}
