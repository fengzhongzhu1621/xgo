package file

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/file"
	"github.com/fengzhongzhu1621/xgo/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
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
func (fs *FileService) ValidateFileType(file io.ReadSeeker) error {
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
	if err := fs.ValidateFileType(src); err != nil {
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

// ZeroCopyDownload 实现零拷贝下载（仅限 Linux 系统）
func (fs *FileService) ZeroCopyDownload(w http.ResponseWriter, filePath string) error {
	fullPath := filepath.Join(fs.storagePath, filePath)

	// 打开文件
	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "无法打开文件", http.StatusInternalServerError)
		return fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "无法获取文件信息", http.StatusInternalServerError)
		return fmt.Errorf("stat file failed: %w", err)
	}

	// 设置 HTTP 头部
	w.Header().
		Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(stat.Name())))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

	// 尝试使用 sendfile 进行零拷贝传输
	if tcpConn, ok := w.(http.Hijacker); ok {
		conn, _, err := tcpConn.Hijack()
		if err != nil {
			http.Error(w, "无法劫持连接", http.StatusInternalServerError)
			return fmt.Errorf("hijack connection failed: %w", err)
		}
		defer conn.Close()

		// 设置 HTTP 响应头
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n" +
			"Content-Disposition: attachment; filename=\"" + filepath.Base(stat.Name()) + "\"\r\n" +
			"Content-Type: application/octet-stream\r\n" +
			"Content-Length: " + fmt.Sprintf("%d", stat.Size()) + "\r\n" +
			"Last-Modified: " + stat.ModTime().UTC().Format(http.TimeFormat) + "\r\n" +
			"\r\n"))
		if err != nil {
			return fmt.Errorf("write response header failed: %w", err)
		}

		// 获取文件描述符
		fileFD := file.Fd()
		connFs, err := conn.(*net.TCPConn).File()
		if err != nil {
			return fmt.Errorf("get net.TCPConn failed: %w", err)
		}
		connFD := connFs.Fd()

		// 计算每次传输的大小（例如 32KB）
		const chunkSize = 32 * 1024
		totalSent := int64(0)

		for totalSent < stat.Size() {
			// 计算本次传输的字节数
			bytesToCopy := int64(chunkSize)
			if totalSent+bytesToCopy > stat.Size() {
				bytesToCopy = stat.Size() - totalSent
			}

			// 使用 sendfile 进行数据传输
			n, err := unix.Sendfile(int(connFD), int(fileFD), nil, int(bytesToCopy))
			if err != nil {
				return fmt.Errorf("sendfile failed: %w", err)
			}

			if n == 0 {
				break // 没有更多数据可发送
			}

			totalSent += int64(n)

			// 如果需要，可以在这里更新进度或其他逻辑
		}

		return nil
	}

	// 如果不支持 Hijacker，则回退到普通复制
	http.ServeContent(w, nil, stat.Name(), stat.ModTime(), file)
	return nil
}

// bufferedCopy 使用缓冲池中的缓冲区进行数据拷贝
func (fs *FileService) bufferedCopy(dst io.Writer, src io.Reader) (int64, error) {
	return file.BufferedCopy(dst, src)
}
