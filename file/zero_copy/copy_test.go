package zero_copy

import (
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"
	"testing"
)

// 文件发送的通用实现。
// 1. 把文件从磁盘读到内核缓冲区。
// 2. 从内核缓冲区复制到程序的内存。
// 3. 从程序内存再复制到网络协议栈。
// 4. 最后发送到客户端。
//
// 优化方案: 采用 io.Copy
func TestIoCopy(t *testing.T) {
	// 打开一个大文件
	// os.File 实现了 io.ReaderFrom，调用了底层的零拷贝机制（比如 sendfile），
	// 直接把文件内容从内核缓冲区送到网络连接。整个过程几乎没有用户空间的缓冲区参与
	file, err := os.Open("big_video.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 连接到服务器
	// net.TCPConn 实现了 io.ReaderFrom，内部可能使用 sendfile 或类似机制
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 零拷贝魔法：直接把文件内容发到网络
	n, err := io.Copy(conn, file)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("成功发送 %d 字节！", n)
}

// CustomSender 自定义网络发送器
type CustomSender struct {
	conn net.Conn
}

// Write 实现 io.Writer 接口
func (s *CustomSender) Write(p []byte) (n int, err error) {
	// 直接调用底层 net.Conn 的 Write 方法，将数据写入网络连接。
	return s.conn.Write(p)
}

// ReadFrom 实现 io.ReaderFrom 接口，利用 sendfile
func (s *CustomSender) ReadFrom(r io.Reader) (n int64, err error) {
	// 仅在 Unix-like 系统上尝试使用 sendfile
	if runtime.GOOS == "windows" {
		// Windows 不支持 sendfile，直接回退
		return io.Copy(s, r)
	}

	// 检查是否为 *os.File
	f, ok := r.(*os.File)
	if !ok {
		// 非文件类型，回退到普通 io.Copy
		return io.Copy(s, r)
	}

	// 获取文件描述符
	fd := int(f.Fd())

	// 获取 TCP 连接的原始套接字
	tcpConn, ok := s.conn.(*net.TCPConn)
	if !ok {
		return 0, io.ErrUnexpectedEOF
	}
	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return 0, err
	}

	// 调用 sendfile
	var totalSent int64
	var offset int64 = 0
	err = rawConn.Control(func(outfd uintptr) {
		for {
			sent, err := syscall.Sendfile(int(outfd), fd, &offset, 4096)
			if err != nil {
				if err == syscall.EAGAIN {
					continue
				}
				break
			}
			totalSent += int64(sent)
			if sent == 0 { // EOF
				break
			}
		}
	})
	if err != nil {
		return totalSent, err
	}
	return totalSent, nil
}

func TestCustomSender(t *testing.T) {
	// 打开文件
	file, err := os.Open("big_video.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 创建自定义发送器
	sender := &CustomSender{conn: conn}

	// 使用 io.Copy 发送文件
	n, err := io.Copy(sender, file)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("成功发送 %d 字节", n)
}
