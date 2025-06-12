package maps

import (
	"fmt"
	"testing"
)

func TestClear(t *testing.T) {
	m := map[string]int{"go": 100, "php": 80}
	fmt.Printf("len=%d\tm=%+v\n", len(m), m) // len=2   m=map[go:100 php:80]
	clear(m)
	fmt.Printf("len=%d\tm=%+v\n", len(m), m) // len=0   m=map[]
}
