package errors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

type MyError struct {
	Code    int
	Message string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func someFunc() error {
	return &MyError{Code: 404, Message: "not found"}
}

func TestCustomError(t *testing.T) {
	err := someFunc()
	if err != nil {
		err = errors.Wrap(err, "additional context") // additional context: Error 404: not found
		fmt.Println(err)
	}
}
