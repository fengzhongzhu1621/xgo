package form_file

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// allowedMimeTypes 允许的 MIME 类型
var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
	"text/plain":      true,
	// 添加其他允许的类型
}

func TestFileService(t *testing.T) {
	router := gin.Default()

	// 初始化 FileService，假设上传文件保存到 "./uploads"
	fileService := NewFileService("./uploads", allowedMimeTypes)

	// 确保上传目录存在
	if _, err := os.Stat(fileService.storagePath); os.IsNotExist(err) {
		os.Mkdir(fileService.storagePath, os.ModePerm)
	}

	// 设置上传路由
	router.POST("/upload", fileService.UploadHandler)

	router.Run(":8080")
}
