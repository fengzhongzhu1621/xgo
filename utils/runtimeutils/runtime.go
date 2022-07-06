package runtimeutils

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// getPackage returns the name of the current package, which makes running this
// test in a fork simpler
func GetPackage() string {
	pc, _, _, _ := runtime.Caller(0)
	fullFuncName := runtime.FuncForPC(pc).Name()
	idx := strings.LastIndex(fullFuncName, ".")
	return fullFuncName[:idx] // trim off function details
}
