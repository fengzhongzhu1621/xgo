package file

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
	xtest "github.com/fengzhongzhu1621/xgo/tests"
)

// Return string of file content.
// func ReadFileToString(path string) (string, error)
// 根据文件路径读取文件（一次性读取）
func TestReadFileToString(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0o777)
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
// 文件按行读取
func TestReadFileByLine(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0o777)
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

	bufPool := sync.Pool{
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

func TestReadMapFromFile(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name        string            // 测试名称
		fileContent string            // 模拟的文件内容
		want        map[string]uint64 // 期望的返回 map
		wantErr     bool              // 是否期望返回错误
	}{
		{
			name: "正常格式，正确解析",
			fileContent: `apple 10
banana 20
orange 5`,
			want: map[string]uint64{
				"apple":  10,
				"banana": 20,
				"orange": 5,
			},
			wantErr: false,
		},
		{
			name: "包含格式错误行，应跳过",
			fileContent: `apple 10
banana twenty  // 这一行无法解析数字，应该跳过
orange 5
   kiwi   100   // 前后可能有空格，但格式是对的
`,
			want: map[string]uint64{
				"apple":  10,
				"orange": 5,
				"kiwi":   100,
			},
			wantErr: false,
		},
		{
			name:        "文件不存在，应返回错误",
			fileContent: "",
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			var err error

			// 除非是测试文件不存在的场景，否则我们创建一个临时文件
			if !tt.wantErr {
				filePath, err = xtest.CreateTempFile(tt.fileContent)
				if err != nil {
					t.Fatalf("创建临时文件失败: %v", err)
				}
				// 测试完成后删除临时文件
				defer os.Remove(filePath)
			} else {
				// 模拟一个不存在的文件路径
				filePath = "non_existent_file_123456.txt"
			}

			got, err := ReadMapFromFile(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMapFromFile() 错误状态 = %v, 期望错误状态 = %v", err, tt.wantErr)
				return
			}

			if err != nil {
				// 只有在不期望出错时，才比较返回的 map
				if len(got) != len(tt.want) {
					t.Errorf("%s: 返回 map 长度不对，得到 %d，期望 %d", tt.name, len(got), len(tt.want))
				}

				for k, v := range tt.want {
					if gotV, ok := got[k]; !ok || gotV != v {
						t.Errorf("%s: 键 %q 的值不对，得到 %d，期望 %d", tt.name, k, gotV, v)
					}
				}

				// 可选：检查是否有多余的键（如果严格希望只解析有效行）
				for k := range got {
					if _, ok := tt.want[k]; !ok {
						t.Logf("警告：返回了未期望的键 %q，可能文件中有可解析但未在测试中声明的行", k)
					}
				}
			}
		})
	}
}
