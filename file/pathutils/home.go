package pathutils

import (
	"os"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// GetHomeDir 获得当前用户的$HOME目录.
func GetHomeDir() string {
	home, _ := homedir.Dir()
	return home
}
