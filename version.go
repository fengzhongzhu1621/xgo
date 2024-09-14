package xgo

import (
	"fmt"
	"regexp"
)

// Version is the current package version.
const (
	Version   = "0.0.1"
	Commit    = "none"
	BuildTime = "unknown"
	GoVersion = "1.23.1"
)

var (
	versionRegex = regexp.MustCompile(`(\d+).(\d+).(\d+)`)
)

// ShowVersion is the default handler which match the --version flag.
func ShowVersion() {
	fmt.Printf("%s", GetVersion())
}

// GetVersion get version message string.
func GetVersion() string {
	version := fmt.Sprintf("Version  :%s\n", Version)
	return version
}

// /////////////////////////////////////////////////////////////////////////
type VersionList [][3]int

func (m VersionList) Len() int {
	return len(m)
}

// Less 版本号比较
func (m VersionList) Less(i, j int) bool {
	for x := range m[i] {
		// 版本号比较
		if m[i][x] == m[j][x] {
			continue
		}
		return m[i][x] < m[j][x]
	}
	return false
}

// Swap 版本号交换
func (m VersionList) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
