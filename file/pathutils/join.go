package pathutils

import (
	"path"
)

// JoinPaths 路径合并.
func JoinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if LastChar(relativePath) == '/' && LastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

