package env

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/system"
	"github.com/stretchr/testify/assert"
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

// func CompareOsEnv(key, comparedEnv string) bool
func TestCompareOsEnv(t *testing.T) {
	err := system.SetOsEnv("foo", "abc")
	if err != nil {
		return
	}

	result := system.CompareOsEnv("foo", "abc")
	assert.Equal(t, true, result)
}
