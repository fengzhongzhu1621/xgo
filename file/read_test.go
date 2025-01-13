package file

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Return string of file content.
// func ReadFileToString(path string) (string, error)
func TestReadFileToString(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0777)
	defer f.Close()

	f.WriteString("hello world")

	content, _ := fileutil.ReadFileToString(fname)

	os.Remove(fname)

	fmt.Println(content)

	// Output:
	// hello world
}

// Read file line by line, and return slice of lines
// func ReadFileByLine(path string)([]string, error)
func TestReadFileByLine(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0777)
	defer f.Close()

	f.WriteString("hello\nworld")

	content, _ := fileutil.ReadFileByLine(fname)

	os.Remove(fname)

	fmt.Println(content)

	// Output:
	// [hello world]
}

// Read File or URL.
// func ReadFile(path string) (reader io.ReadCloser, closeFn func(), err error)
func TestReadFile(t *testing.T) {
	reader, fn, err := fileutil.ReadFile("https://httpbin.org/robots.txt")
	if err != nil {
		return
	}
	defer fn()

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

// TestChunkRead reads a block from the file at the specified offset and returns all lines within the block.
// func ChunkRead(file *os.File, offset int64, size int, bufPool *sync.Pool) ([]string, error)
func TestChunkRead(t *testing.T) {
	fpath := "./test.csv"
	fileutil.CreateFile(fpath)

	data := [][]string{
		{"Lili", "22", "female"},
		{"Jim", "21", "male"},
	}
	err := fileutil.WriteCsvFile(fpath, data, false)

	if err != nil {
		return
	}

	const mb = 1024 * 1024
	const defaultChunkSizeMB = 100

	f, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer f.Close()

	var bufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, defaultChunkSizeMB*mb)
		},
	}

	lines, err := fileutil.ChunkRead(f, 0, 100, &bufPool)
	if err != nil {
		return
	}

	fmt.Println(lines[0])
	fmt.Println(lines[1])

	// Output:
	// Lili,22,female
	// Jim,21,male
}

// Reads the file in parallel and send each chunk of lines to the specified channel.
// filePath: file path.
// chunkSizeMB: The size of the block (in MB, the default is 100MB when set to 0). Setting it too large will be detrimental. Adjust it as appropriate.
// maxGoroutine: The number of concurrent read chunks, the number of CPU cores used when set to 0.
// linesCh: The channel used to receive the returned results.
// func ParallelChunkRead(filePath string, linesCh chan<- []string, chunkSizeMB, maxGoroutine int) error
func TestParallelChunkRead(t *testing.T) {
	fpath := "./test.csv"
	fileutil.CreateFile(fpath)

	data := [][]string{
		{"Lili", "22", "female"},
		{"Jim", "21", "male"},
	}
	err := fileutil.WriteCsvFile(fpath, data, false)

	if err != nil {
		return
	}

	const mb = 1024 * 1024
	const defaultChunkSizeMB = 100 // 默认值

	numParsers := runtime.NumCPU()

	linesCh := make(chan []string, numParsers)

	go fileutil.ParallelChunkRead(fpath, linesCh, defaultChunkSizeMB, numParsers)

	var totalLines int
	for lines := range linesCh {
		totalLines += len(lines)

		for _, line := range lines {
			fmt.Println(line)
		}
	}

	fmt.Println(totalLines)

	// Output:
	// Lili,22,female
	// Jim,21,male
	// 2

}
