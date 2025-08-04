package uploader

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
)

// UploadChunks 并发上传所有分片
func UploadChunks(config Config, chunks [][]byte) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(chunks))                // 缓冲 channel 避免 goroutine 阻塞
	semaphore := make(chan struct{}, config.MaxConcurrency) // 控制并发数

	// 遍历每个分片
	for i, chunk := range chunks {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(partNumber int, filePart []byte) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			// 重试逻辑
			var lastErr error
			for retry := 0; retry < config.MaxRetries; retry++ {
				err := uploadChunk(
					config.UploadURL,
					config.FilePath,
					filePart,
					partNumber,
					&wg,
					errChan,
				)
				if err == nil {
					return // 上传成功
				}
				lastErr = err
			}
			errChan <- fmt.Errorf("part %d failed after retries: %v", partNumber, lastErr)
		}(i+1, chunk) // partNumber 从 1 开始
	}

	wg.Wait()
	close(errChan)

	// 检查是否有错误
	var hasError bool
	for err := range errChan {
		fmt.Println("Upload error:", err)
		hasError = true
	}

	if hasError {
		return ErrUploadFailed
	}
	return nil
}

// uploadChunk 上传单个分片（和你之前提供的代码几乎一致，稍作调整）
// UploadChunk 上传块的函数
// url: 上传的目标服务器地址。
// filename: 原始文件名（用于构造分块文件的名称）。
// filePart: 当前要上传的文件块数据（[]byte 类型）。
// partNumber: 当前块的编号（用于标识是第几块）。
// wg *sync.WaitGroup: 用于等待所有分块上传完成的同步机制。
// errChan chan error: 用于传递上传过程中发生的错误。
func uploadChunk(
	url, filename string,
	filePart []byte,
	partNumber int,
	wg *sync.WaitGroup,
	errChan chan error,
) error {
	defer wg.Done()

	// 创建 multipart 表单数据
	// 使用 multipart.NewWriter 创建一个 multipart 写入器
	// 创建一个表单字段 "file"，字段值是类似 filename.part1 的形式
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fmt.Sprintf("%s.part%d", filename, partNumber))
	if err != nil {
		errChan <- err
		return err
	}

	// 写入文件块数据，将当前块的二进制数据写入到 multipart 表单中
	if _, err := part.Write(filePart); err != nil {
		errChan <- err
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		errChan <- err
		return err
	}

	// 构造 HTTP 上传请求
	// 设置 Content-Type 为 multipart 的格式（包含 boundary 等信息）。
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	// 发送请求并处理响应
	resp, err := client.Do(req)
	if err != nil {
		errChan <- err
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := string(bodyBytes)
		errChan <- fmt.Errorf("failed to upload part %d, status: %s, response: %s", partNumber, resp.Status, errMsg)
		return fmt.Errorf("server returned non-200 for part %d: %s", partNumber, resp.Status)
	}

	return nil
}
