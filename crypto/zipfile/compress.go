package zipfile

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	dkignore "github.com/codeskyblue/dockerignore"
)

const YAMLCONF = ".ghs.yml"

// 压缩目录，写入到 http 响应中，返回给前端压缩后的.zip文件
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

func ExtractFromZip(zipFile, path string, w io.Writer) (err error) {
	// 用于打开 ZIP 文件并返回一个 *zip.Reader 类型的值
	cf, err := zip.OpenReader(zipFile)
	if err != nil {
		return
	}
	defer cf.Close()

	// 匹配规则内容：接受一个字符串参数并返回一个具有指定字符串内容的缓冲区（*bytes.Buffer 类型）
	rd := io.NopCloser(bytes.NewBufferString(path))
	// 读取匹配规则
	patterns, err := dkignore.ReadIgnore(rd)
	if err != nil {
		return
	}

	// 遍历压缩包中的文件
	for _, file := range cf.File {
		// 判断文件是否匹配
		matched, _ := dkignore.Matches(file.Name, patterns)
		if !matched {
			continue
		}
		// 读取文件
		rc, er := file.Open()
		if er != nil {
			err = er
			return
		}
		defer rc.Close()
		// 复制文件内容
		_, err = io.Copy(w, rc)
		if err != nil {
			return
		}
		return
	}

	return fmt.Errorf("file %s not found", strconv.Quote(path))
}
