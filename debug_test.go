package xgo

import (
	"runtime"
	"testing"
)

func TestDebugPrintWARNINGNew(t *testing.T) {
	debugPrintWARNINGNew()
}

func TestDebugPrintWARNINGDefault(t *testing.T) {
	t.Log("git version = " + runtime.Version())
	debugPrintWARNINGDefault()
}
