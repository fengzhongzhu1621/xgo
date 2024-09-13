package ssl

import (
	"errors"
	"strconv"
)

type SSHVersion struct {
	Major int
	Minor int
}

// ExtractOpenSSHVersion 提取 OpenSSH 版本号
func ExtractOpenSSHVersion(input string) ([]string, error) {
	// 找到了所有的匹配项及其子匹配项，并将它们存储在一个二维切片中
	matches := OpensslVersionRegex.FindAllStringSubmatch(input, -1)

	versions := make([]string, 0, len(matches))

	// 提取 OpenSSH 字符后面的版本号
	for _, match := range matches {
		versions = append(versions, match[1])
	}

	return versions, nil
}

// ExtractOpenSSHVersionV2 提取 OpenSSH 版本号（包含主版本号和子版本号）
func ExtractOpenSSHVersionWithMajorAndMinor(input string) (*SSHVersion, error) {
	// 找到了所有的匹配项及其子匹配项，并将它们存储在一个二维切片中
	matches := OpensslVersionRegexV2.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		major, _ := strconv.Atoi(match[1])
		minor, _ := strconv.Atoi(match[2])
		return &SSHVersion{
			Major: major,
			Minor: minor}, nil
	}

	return nil, errors.New("version format error")
}
