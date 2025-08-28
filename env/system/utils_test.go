package system

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/system"
)

// func GetOsEnv(key string) string
// func SetOsEnv(key, value string) error
// func RemoveOsEnv(key string) error
func TestGetOsEnv(t *testing.T) {
	err := system.SetOsEnv("foo", "abc")
	result := system.GetOsEnv("foo")
	// <nil>
	fmt.Println(err)
	// abc
	fmt.Println(result)

	err2 := system.RemoveOsEnv("foo")
	result2 := system.GetOsEnv("foo")
	// <nil>
	fmt.Println(err2)
	//
	fmt.Println(result2)
}
