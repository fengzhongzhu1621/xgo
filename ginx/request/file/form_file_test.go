package file

import (
	"fmt"
	"net/http"
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
