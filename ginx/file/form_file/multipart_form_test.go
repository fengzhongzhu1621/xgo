package form_file

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMultipartForm(t *testing.T) {
	engine := gin.Default()

	// 设置最大 multipart 内存为 8MB（超过部分会存到临时文件）
	engine.MaxMultipartMemory = 8 << 20 // 8MB

	engine.POST("/uploadMul", func(c *gin.Context) {
		// 1. 解析 multipart 表单
		form, err := c.MultipartForm()
		if err != nil {
			log.Println("解析 multipart 表单失败:", err)
			c.String(http.StatusBadRequest, "文件上传失败: %v", err)
			return
		}

		// 2. 获取上传的文件列表
		files := form.File["upload"]
		if len(files) == 0 {
			c.String(http.StatusBadRequest, "没有上传任何文件")
			return
		}

		// 3. 确保上传目录存在
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0o755); err != nil {
			log.Println("创建上传目录失败:", err)
			c.String(http.StatusInternalServerError, "服务器内部错误")
			return
		}

		// 4. 遍历并保存每个文件
		// // 在保存文件前添加检查
		const maxFileSize = 4 << 20 // 4MB
		allowedExtensions := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}

		var savedFiles []string
		for _, file := range files {
			// 安全处理文件名（防止路径遍历攻击）
			safeFilename := filepath.Base(file.Filename)
			dst := filepath.Join(uploadDir, safeFilename)

			ext := filepath.Ext(safeFilename)

			// 检查文件扩展名
			if !allowedExtensions[ext] {
				log.Printf("拒绝文件 %s: 不支持的文件类型", file.Filename)
				continue
			}

			// 检查文件大小
			if file.Size > maxFileSize {
				log.Printf("拒绝文件 %s: 文件大小超过限制 (%d > %d)", file.Filename, file.Size, maxFileSize)
				continue
			}

			// 打开上传的文件
			srcFile, err := file.Open()
			if err != nil {
				log.Printf("打开文件 %s 失败: %v", file.Filename, err)
				continue // 跳过当前文件，继续处理下一个
			}

			// 创建目标文件
			dstFile, err := os.Create(dst)
			if err != nil {
				log.Printf("创建文件 %s 失败: %v", dst, err)
				srcFile.Close()
				continue
			}

			// 复制文件内容
			if _, err := io.Copy(dstFile, srcFile); err != nil {
				log.Printf("保存文件 %s 失败: %v", dst, err)
				srcFile.Close()
				dstFile.Close()
				continue
			}

			// 关闭文件句柄
			srcFile.Close()
			dstFile.Close()

			savedFiles = append(savedFiles, safeFilename)
		}

		// 5. 返回结果
		if len(savedFiles) == 0 {
			c.String(http.StatusInternalServerError, "所有文件上传失败")
			return
		}

		c.String(
			http.StatusOK,
			fmt.Sprintf("成功上传 %d/%d 个文件: %v", len(savedFiles), len(files), savedFiles),
		)
	})

	engine.Run(":8080")
}
