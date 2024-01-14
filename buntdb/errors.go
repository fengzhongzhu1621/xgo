package buntdb

import "fmt"

func panicErr(err error) error {
	panic(fmt.Errorf("buntdb: %w", err))
}
