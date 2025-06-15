package uploader

import "errors"

// 配置参数
type Config struct {
	FilePath       string // 待上传的文件路径
	ChunkSize      int64  // 分片大小（字节），如 5MB = 5 * 1024 * 1024
	UploadURL      string // 上传分片的 URL
	MergeURL       string // 合并分片的 URL（可选）
	MaxConcurrency int    // 最大并发上传数
	MaxRetries     int    // 单个分片最大重试次数
}

// 自定义错误
var (
	ErrFileTooSmall = errors.New("file is too small to split")
	ErrUploadFailed = errors.New("some chunks failed to upload")
)
