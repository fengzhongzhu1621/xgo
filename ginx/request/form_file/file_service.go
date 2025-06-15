package form_file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/fengzhongzhu1621/xgo/validator"
)

// FileService 结构体
type FileService struct {
	storagePath      string
	allowedMimeTypes map[string]bool
}

// NewFileService 创建一个新的 FileService 实例
func NewFileService(storagePath string, allowedMimeTypes map[string]bool) *FileService {
	return &FileService{storagePath: storagePath, allowedMimeTypes: allowedMimeTypes}
}

// validateFileType 使用 filetype 库检测 MIME 类型
func (fs *FileService) validateFileType(file io.ReadSeeker) error {
	err := validator.ValidateFileType(file, fs.allowedMimeTypes)
	if err != nil {
		return err
	}

	return nil
}

// UploadHandler 处理文件上传
func (fs *FileService) UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取文件"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件"})
		return
	}
	defer src.Close()

	// 验证文件类型
	if err := fs.validateFileType(src); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存文件到指定路径
	dst, err := os.Create(filepath.Join(fs.storagePath, file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存文件"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法复制文件内容"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("文件 '%s' 上传成功", file.Filename)})
}
