package pathutils

import "os"

// GetWd 获得应用程序当前路径.
func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}
