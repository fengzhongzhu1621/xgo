package file

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func TestUploadHandler(t *testing.T) {
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

func TestZeroCopyDownload(t *testing.T) {
	fs := &FileService{
		storagePath: "./uploads", // 替换为你的存储路径
	}

	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		// 解析文件路径，确保安全性（防止目录遍历）
		path := r.URL.Path[len("/download/"):]
		if path == "" {
			http.Error(w, "文件路径不能为空", http.StatusBadRequest)
			return
		}
		fullPath := filepath.Join(fs.storagePath, filepath.Clean(path))

		// 检查文件是否存在及是否为文件
		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			http.Error(w, "文件不存在", http.StatusNotFound)
			return
		}
		if !info.Mode().IsRegular() {
			http.Error(w, "不是一个常规文件", http.StatusBadRequest)
			return
		}

		// 执行零拷贝下载
		if err := fs.ZeroCopyDownload(w, path); err != nil {
			http.Error(w, "下载失败", http.StatusInternalServerError)
			fmt.Println("ZeroCopyDownload error:", err)
		}
	})

	fmt.Println("服务器启动，监听端口 :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
