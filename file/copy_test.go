package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/fengzhongzhu1621/xgo/validator"
	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	{
		var dst string = "./file_test2.go"
		err := Copy("./file_test.go", dst)
		fmt.Println(err)
		assert.Equal(t, validator.IsFileOrDirExists(dst), true)
		err = os.Remove(dst)
		assert.NoError(t, err)
	}

	{
		// TestCopyFile Copy src file to dest file. If dest file exist will overwrite it.
		// func CopyFile(srcPath string, dstPath string) error
		srcFile := "./text.txt"
		fileutil.CreateFile(srcFile)

		copyFile := "./text_copy.txt"
		fileutil.CopyFile(srcFile, copyFile)

		file, err := os.Open(copyFile)
		if err != nil {
			return
		}
		result1 := fileutil.IsExist(copyFile)
		result2 := file.Name()

		os.Remove(srcFile)
		os.Remove(copyFile)

		fmt.Println(result1)
		fmt.Println(result2)

		// Output:
		// true
		// ./text_copy.txt
	}
}

// TestCopyDir 将src目录复制到dst目录，它将递归地复制所有文件和目录。访问权限将与源目录相同。如果dstPath已存在，则会返回一个错误。
// func CopyDir(srcPath string, dstPath string) error
func TestCopyDir(t *testing.T) {
	{
		// 定义源目录和目标目录
		srcDir := "source_dir"
		destDir := "destination_dir"

		// 使用 os.CopyFS 复制目录
		err := os.CopyFS(destDir, os.DirFS(srcDir))
		if err != nil {
			log.Fatalf("Failed to copy directory: %v", err)
		}
		log.Println("Directory copied successfully!")
	}

	{
		pwd, _ := os.Getwd()
		srcPath := filepath.Join(pwd, "test_src")
		if !fileutil.IsExist(srcPath) {
			os.RemoveAll(srcPath)
		}

		fileutil.CreateDir(srcPath)
		fileutil.CreateFile(filepath.Join(srcPath, "test.txt01"))
		fileutil.CreateFile(filepath.Join(srcPath, "test.txt02"))
		fileutil.CreateDir(filepath.Join(srcPath, "test_dir01"))
		fileutil.CreateDir(filepath.Join(srcPath, "test_dir02"))
		fileutil.CreateDir(filepath.Join(srcPath, "test_dir02", "text.txt03"))

		destPath := filepath.Join(pwd, "test_dest")

		err := fileutil.CopyDir(srcPath, destPath)
		if err != nil {
			fmt.Printf("copy dir error: %v", err)
			return
		}

		// check dest path exist
		if !fileutil.IsExist(destPath) {
			fmt.Printf("copy dir error: dest path not exist")
			return
		}
		filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
			destp := strings.Replace(path, srcPath, destPath, 1)
			if !fileutil.IsExist(destp) {
				fmt.Printf("copy dir error: %s not exist", destp)
				return nil
			}
			return nil
		})

		os.RemoveAll(destPath)
	}
}
