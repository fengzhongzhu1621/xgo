package zipfile

import (
	"archive/zip"
	"net/http"
	"os"
	"path/filepath"
)

const YAMLCONF = ".ghs.yml"

// CompressToZip 压缩目录，写入到 http 响应中，返回给前端压缩后的.zip文件
func CompressToZip(w http.ResponseWriter, rootDir string) {
	// 用于清理目录路径，它将路径表示为规范形式，并消除任何多余的斜杠和.、..元素。
	rootDir = filepath.Clean(rootDir)
	// 获得压缩包的名称，用于获取路径的最后一个元素。这个元素可以是目录名或文件名。
	zipFileName := filepath.Base(rootDir) + ".zip"

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+zipFileName+`"`)

	// 创建 zip 对象
	zw := &Zip{Writer: zip.NewWriter(w)}
	defer zw.Close()

	// 遍历文件夹
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		// 获得文件的相对路径
		zipPath := path[len(rootDir):]
		// 根据规则忽略指定文件
		if info.Name() == YAMLCONF { // ignore .ghs.yml for security
			return nil
		}
		// 将文件添加到压缩包
		return zw.Add(zipPath, path)
	})
}
