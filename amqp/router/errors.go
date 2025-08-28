package router

import (
	"fmt"
)

// DuplicateHandlerNameError is sent in a panic when you try to add a second handler with the same name.
type DuplicateHandlerNameError struct {
	HandlerName string
}

func (d DuplicateHandlerNameError) Error() string {
	return fmt.Sprintf("handler with name %s already exists", d.HandlerName)
}
