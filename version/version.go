package version

import (
	"fmt"
	"regexp"
	"runtime"
)

// Version is the current package version.
var (
	// AppVersion 版本号
	AppVersion = "--"
	// GitCommit CommitID
	GitCommit = "--"
	// BuildTime 二进制构建时间
	BuildTime = "--"
	// 模板的版本号
	TemplateVersion = "0.0.1"
	GoVersion       = runtime.Version()
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
	version := fmt.Sprintf("Version  :%s\n", AppVersion)
	return version
}

// /////////////////////////////////////////////////////////////////////////
func Version() string {
	return fmt.Sprintf(
		"\nVersion  : %s\nGitCommit: %s\nBuildTime: %s\nTemplateVersion: %s\nGoVersion: %s\n",
		AppVersion, GitCommit, BuildTime, TemplateVersion, GoVersion,
	)
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
