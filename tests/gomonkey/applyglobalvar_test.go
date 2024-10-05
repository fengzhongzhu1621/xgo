package gomonkey

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

var num = 10

func TestGlobalVar(t *testing.T) {
	patches := gomonkey.ApplyGlobalVar(&num, 12)
	defer patches.Reset()

	if num != 12 {
		t.Errorf("expected %v, got %v", 12, num)
	}
}
