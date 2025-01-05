package shell

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/system"
)

func TestIsLinux(t *testing.T) {
	isOsLinux := system.IsLinux()
	fmt.Println(isOsLinux)
}

func TestIsMac(t *testing.T) {
	isOsMac := system.IsMac()
	fmt.Println(isOsMac)
}

// GetOsBits  Get current os bits, 32bit or 64bit. return 32 or 64.
// func GetOsBits() int
func TestGetOsBits(t *testing.T) {
	osBit := system.GetOsBits()
	fmt.Println(osBit) // 32 or 64
}
