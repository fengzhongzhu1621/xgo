package form_file

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestFormUpload(t *testing.T) {
	// 创建默认的Gin路由引擎
	router := gin.Default()

	// 设置multipart表单的内存限制为8MB（默认是32MB）
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 处理文件上传的POST请求
	router.POST("/upload", func(c *gin.Context) {
		// 从表单中获取上传的文件
		file, _ := c.FormFile("file")

		// 获取文件名（不含路径）
		filename := filepath.Base(file.Filename)

		// 尝试将上传的文件保存到服务器
		if err := c.SaveUploadedFile(file, filename); err != nil {
			// 如果保存失败，返回400错误和错误信息
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		// 文件保存成功，返回200状态码和成功信息
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", filename))
	})

	// 启动HTTP服务器，监听8080端口
	router.Run(":8080")
}

// 测试正常上传
// curl -X POST -F "file=@test.jpg" http://localhost:8080/upload
//
// 测试文件过大
// # 生成一个 5MB 的测试文件
// dd if=/dev/urandom of=large_file.jpg bs=1M count=5
// curl -X POST -F "file=@large_file.jpg" http://localhost:8080/upload
//
// 测试非法文件类型
// # 创建一个文本文件
// echo "test" > test.txt
// curl -X POST -F "file=@test.txt" http://localhost:8080/upload

func TestSingleFileUpload(t *testing.T) {
	engine := gin.Default()

	// 设置最大 multipart 内存为 8MB（超过部分会存到临时文件）
	engine.MaxMultipartMemory = 8 << 20 // 8MB

	engine.POST("/upload", func(c *gin.Context) {
		// 1. 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			log.Println("获取文件失败:", err)
			c.String(http.StatusBadRequest, "文件上传失败: %v", err)
			return
		}

		// 2. 检查文件大小（限制为 4MB）
		const maxFileSize = 4 << 20 // 4MB
		if file.Size > maxFileSize {
			c.String(http.StatusBadRequest, "文件大小不能超过 4MB")
			return
		}

		// 3. 检查文件类型（通过扩展名验证）
		allowedExtensions := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}

		ext := filepath.Ext(file.Filename)
		if !allowedExtensions[ext] {
			c.String(http.StatusBadRequest, "只允许上传 JPG/JPEG/PNG 格式图片")
			return
		}

		// 4. 直接保存文件（不提前读取内容）
		dst := "./uploads/" + file.Filename
		if err := os.MkdirAll("./uploads", 0o755); err != nil {
			log.Println("创建上传目录失败:", err)
			c.String(http.StatusInternalServerError, "服务器内部错误")
			return
		}
		if err := c.SaveUploadedFile(file, dst); err != nil {
			log.Println("保存文件失败:", err)
			c.String(http.StatusInternalServerError, "文件保存失败: %v", err)
			return
		}

		c.String(http.StatusOK, "'%s' 上传成功!", file.Filename)
	})

	engine.Run(":8080")
}
