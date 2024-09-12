package xgo

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

// /////////////////////////////////////////////////////////////////////////
type VersionLog struct {
	Version     string `json:"version"`
	ReleaseDate string `json:"release_at"`
	Content     string `json:"content"`
}

// GetLatestVersion 获得多个版本字符串中的最新版本
func GetLatestVersion(vers []string) string {
	var (
		versionListVal [][3]int
		defaultVersion string
		err            error
	)

	for _, v := range vers {
		if v == "" {
			continue
		}
		// 正则表达式匹配版本号
		matchs := versionRegex.FindStringSubmatch(v)
		if len(matchs) > 0 && len(matchs) == 4 {
			defaultVersion = fmt.Sprintf("V%s", matchs[0])
			var (
				a, b, c int
			)
			if a, err = strconv.Atoi(matchs[1]); err != nil {
				continue
			}
			if b, err = strconv.Atoi(matchs[2]); err != nil {
				continue
			}
			if c, err = strconv.Atoi(matchs[3]); err != nil {
				continue
			}
			versionListVal = append(versionListVal, [3]int{a, b, c})
		}
	}

	if len(versionListVal) > 0 {
		versionSort := VersionList(versionListVal)
		sort.Sort(&versionSort)
		// 获得最新的版本号
		final := versionSort[len(versionSort)-1]

		return fmt.Sprintf("V%d.%d.%d", final[0], final[1], final[2])
	}

	return defaultVersion
}

// ReadFileContent 读取文件的内容
func readFileContent(filepath string) (string, error) {
	// 用于读取指定文件的全部内容，并将其作为一个字节切片返回。如果读取过程中发生错误，函数会返回一个非空的错误值。
	// os.ReadFile 函数会一次性读取整个文件内容到内存中，因此对于非常大的文件，可能会导致内存不足的问题。
	// 在这种情况下，建议使用 os.Open 和 io.Reader 接口逐块读取文件内容。
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ListChangelogs(rootDir string, language string, version string) (string, []*VersionLog, error) {
	var (
		err                error
		files              []string
		versionLogs        []*VersionLog
		content, latestVer string
	)

	// 获得所有后缀为 .md 的文件名
	// 用于遍历指定目录下的所有文件和子目录。它会递归地遍历目录树，并对每个文件和目录调用指定的回调函数。
	// 如果在遍历过程中发生错误，返回非空的错误值；否则返回 nil
	//
	// type WalkFunc func(path string, info os.FileInfo, err error) error
	// - path：当前遍历到的文件或目录的完整路径。
	// - info：当前遍历到的文件或目录的 os.FileInfo 信息。
	// - err：在遍历过程中遇到的错误，如果为 nil，则表示没有错误发生
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		// 忽略目录
		if info.IsDir() {
			return nil
		}

		// 忽略非 markdown 文件
		if filepath.Ext(path) != ".md" {
			return nil
		}
		files = append(files, info.Name())
		return nil
	})
	if err != nil {
		return "", nil, err
	}

	// 从版本文件名中获得最新的版本号
	latestVer = GetLatestVersion(files)

	for _, filename := range files {
		// 判断路径是否包含指定语言
		// V1.0.0_20230227_en
		// V1.0.0_20230227_zh-hans
		if language != "" {
			if !strings.Contains(filename, language) {
				continue
			}
		}
		filename_arr := strings.Split(filename, "_")
		if len(filename_arr) < 3 {
			continue
		}

		// 读取 md 文件
		if version != "" && strings.Contains(filename, version) && strings.Contains(filename, language) {
			content, err = readFileContent(rootDir + filename)
			if err != nil {
				continue
			}
		}

		versionLogs = append(versionLogs, &VersionLog{
			Version:     filename_arr[0],
			ReleaseDate: filename_arr[1],
			Content:     content,
		})
	}

	return latestVer, versionLogs, nil
}
