package zipfile

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Create a zip file of fpath, fpath could be a file or a directory.
// func Zip(fpath string, destPath string) error
func TestZip(t *testing.T) {
	srcFile := "./test.txt"
	fileutil.CreateFile(srcFile)

	zipFile := "./test.zip"
	fileutil.Zip(srcFile, zipFile)

	result := fileutil.IsExist(zipFile)

	os.Remove(srcFile)
	os.Remove(zipFile)

	fmt.Println(result)

	// Output:
	// true
}

// Append a single file or directory by fpath to an existing zip file.
// func ZipAppendEntry(fpath string, destPath string) error
func TestZipAppendEntry(t *testing.T) {
	zipFile := "./test.zip"
	fileutil.CopyFile("./testdata/file.go.zip", zipFile)

	fileutil.ZipAppendEntry("./testdata", zipFile)

	unZipPath := "./unzip"
	fileutil.UnZip(zipFile, unZipPath)

	fmt.Println(fileutil.IsExist("./unzip/file.go"))
	fmt.Println(fileutil.IsExist("./unzip/testdata/file.go.zip"))
	fmt.Println(fileutil.IsExist("./unzip/testdata/test.csv"))
	fmt.Println(fileutil.IsExist("./unzip/testdata/test.txt"))

	os.Remove(zipFile)
	os.RemoveAll(unZipPath)

	// Output:
	// true
	// true
	// true
	// true
}

// Unzip the file and save it to dest path.
// func UnZip(zipFile string, destPath string) error
func TestUnZip(t *testing.T) {
	srcFile := "./test.txt"
	fileutil.CreateFile(srcFile)

	zipFile := "./test.zip"
	fileutil.Zip(srcFile, zipFile)

	unZipPath := "./unzip"
	fileutil.UnZip(zipFile, unZipPath)

	result := fileutil.IsExist("./unzip/test.txt")

	os.Remove(srcFile)
	os.Remove(zipFile)
	os.RemoveAll(unZipPath)

	fmt.Println(result)

	// Output:
	// true
}
