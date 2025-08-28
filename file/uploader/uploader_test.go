package uploader

import (
	"fmt"
	"log"
	"testing"
)

func TestUploader(t *testing.T) {
	config := Config{
		FilePath:       "test.zip",                     // 待上传的文件
		ChunkSize:      5 * 1024 * 1024,                // 5MB
		UploadURL:      "http://localhost:8080/upload", // 分片上传 URL
		MergeURL:       "http://localhost:8080/merge",  // 合并 URL（可选）
		MaxConcurrency: 3,                              // 最大并发数
		MaxRetries:     3,                              // 单个分片最大重试次数
	}

	// 1. 分片文件
	chunks, err := SplitFile(config.FilePath, config.ChunkSize)
	if err != nil {
		log.Fatal("SplitFile failed:", err)
	}

	fmt.Printf("File split into %d chunks\n", len(chunks))

	// 2. 并发上传分片
	if err := UploadChunks(config, chunks); err != nil {
		log.Fatal("UploadChunks failed:", err)
	}

	fmt.Println("\nAll chunks uploaded successfully!")
}
