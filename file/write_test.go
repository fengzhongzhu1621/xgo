package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// TestWriteBytesToFile Writes bytes to target file.
// func WriteBytesToFile(filepath string, content []byte) error
func TestWriteBytesToFile(t *testing.T) {
	filepath := "./bytes.txt"

	file, err := os.Create(filepath)
	if err != nil {
		return
	}

	defer file.Close()

	err = fileutil.WriteBytesToFile(filepath, []byte("hello"))
	if err != nil {
		return
	}

	content, err := fileutil.ReadFileToString(filepath)
	if err != nil {
		return
	}

	os.Remove(filepath)

	fmt.Println(content)

	// Output:
	// hello
}

// TestWriteStringToFile Writes string to target file.
// func WriteStringToFile(filepath string, content string, append bool) error
func TestWriteStringToFile(t *testing.T) {
	filepath := "./test.txt"

	file, err := os.Create(filepath)
	if err != nil {
		return
	}

	defer file.Close()

	err = fileutil.WriteStringToFile(filepath, "hello", true)
	if err != nil {
		return
	}

	content, err := fileutil.ReadFileToString(filepath)
	if err != nil {
		return
	}

	os.Remove(filepath)

	fmt.Println(content)

	// Output:
	// hello
}
