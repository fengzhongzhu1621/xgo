package uploader

import (
	"io"
	"os"
)

// SplitFile 将文件按 ChunkSize 分片，返回分片数据的切片
func SplitFile(filePath string, chunkSize int64) ([][]byte, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获得文件大小，并判断分片的大小
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()
	if fileSize <= chunkSize {
		return nil, ErrFileTooSmall
	}

	var chunks [][]byte
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		// 如果是最后一次读取，可能不足 chunkSize
		if err == io.EOF {
			chunk := make([]byte, n)
			copy(chunk, buffer[:n])
			chunks = append(chunks, chunk)
		} else {
			chunks = append(chunks, buffer[:n])
		}
	}

	return chunks, nil
}
