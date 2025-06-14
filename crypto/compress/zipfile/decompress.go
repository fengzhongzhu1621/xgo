package zipfile

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	dkignore "github.com/codeskyblue/dockerignore"
	"github.com/fengzhongzhu1621/xgo/file"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// ExtractFromZip 从压缩包中解压指定路径的文件
func ExtractFromZip(zipFile, path string, w io.Writer) (err error) {
	// 用于打开 ZIP 文件并返回一个 *zip.Reader 类型的值
	cf, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil
	}
	defer cf.Close()

	// 匹配规则内容：接受一个字符串参数并返回一个具有指定字符串内容的缓冲区（*bytes.Buffer 类型）
	rd := io.NopCloser(bytes.NewBufferString(path))
	// 读取匹配规则
	patterns, err := dkignore.ReadIgnore(rd)
	if err != nil {
		return nil
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

		// 只处理一个文件
		return
	}

	return fmt.Errorf("file %s not found", strconv.Quote(path))
}

func UnzipFile(filename, dest string) error {
	// 用于打开 ZIP 文件并返回一个 *zip.Reader 类型的值
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zr.Close()

	// 获取文件路径中的目录部分
	if dest == "" {
		dest = filepath.Dir(filename)
	}

	for _, f := range zr.File {
		// 读取文件
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// ignore .ghs.yml
		filename := file.SanitizedName(f.Name)
		if filepath.Base(filename) == ".ghs.yml" {
			continue
		}
		// 获得文件的路径
		fpath := filepath.Join(dest, filename)

		// filename maybe GBK or UTF-8
		// Ref: https://studygolang.com/articles/3114
		/*
			在新的 zip规范 中,
			已经对UTF8编码的文件名提供了支持.

			File:    APPNOTE.TXT - .ZIP File Format Specification
			Version: 6.3.3

			4.4.4 general purpose bit flag: (2 bytes)

			Bit 11: Language encoding flag (EFS).  If this bit is set,
			    the filename and comment fields for this file
			    MUST be encoded using UTF-8. (see APPENDIX D)
			具体来说, 在每个文件的头信息的Flags字段的11bit位.
			如果该bit位为0则表用本地编码(本地编码是GBK吗?), 如果是1则表示用UTF8编码.
		*/
		if f.Flags&(1<<11) == 0 { // GBK
			tr := simplifiedchinese.GB18030.NewDecoder()
			fpathUtf8, err := tr.String(fpath)
			if err == nil {
				fpath = fpathUtf8
			}
		}

		// 创建目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)

		// 创建并打开，并清空文件内容
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		// 复制文件
		_, err = io.Copy(outFile, rc)
		outFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
