package xgo

import "fmt"

// Version is the current package version.
const version = "0.0.1"

// ShowVersion is the default handler which match the --version flag.
func ShowVersion() {
	fmt.Printf("%s", GetVersion())
}

// GetVersion get version message string.
func GetVersion() string {
	version := fmt.Sprintf("Version  :%s\n", version)
	return version
}
