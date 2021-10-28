package pipe

import (
	"bufio"
	"io"
	"syscall"
	"time"

	. "xgo/command"
	. "xgo/log"
	. "xgo/utils/bytesutils"
)

type ProcessEndpoint struct {
	process   *LaunchedProcess // 命令加载进程类
	closetime time.Duration
	output    chan []byte // 标准输出读取的数据发给这个管道
	log       *LogScope
	bin       bool // 传递数据的格式：字符串或二进制
}

// 构造函数
func NewProcessEndpoint(process *LaunchedProcess, bin bool, log *LogScope) *ProcessEndpoint {
	return &ProcessEndpoint{
		process: process,
		output:  make(chan []byte),
		log:     log,
		bin:     bin,
	}
}

func (pe *ProcessEndpoint) Terminate() {
	// 等待命令执行完成
	terminated := make(chan struct{})
	go func() {
		pe.process.Cmd.Wait()
		terminated <- struct{}{}
	}()

	// for some processes this is enough to finish them...
	// 关闭输入，命令停止传入数据
	pe.process.Stdin.Close()

	// a bit verbose to create good debugging trail
	// 命令执行完成直接退出
	select {
	case <-terminated:
		pe.log.Debug("process", "Process %v terminated after stdin was closed", pe.process.Cmd.Process.Pid)
		return // means process finished
	case <-time.After(100*time.Millisecond + pe.closetime):
	}

	// 超时后：发送SIGINT信号 CTRL+C
	err := pe.process.Cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		// process is done without this, great!
		pe.log.Error("process", "SIGINT unsuccessful to %v: %s", pe.process.Cmd.Process.Pid, err)
	}

	// 超时后：发送SIGTERM信号 KILL
	select {
	case <-terminated:
		pe.log.Debug("process", "Process %v terminated after SIGINT", pe.process.Cmd.Process.Pid)
		return // means process finished
	case <-time.After(250*time.Millisecond + pe.closetime):
	}

	err = pe.process.Cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		// process is done without this, great!
		pe.log.Error("process", "SIGTERM unsuccessful to %v: %s", pe.process.Cmd.Process.Pid, err)
	}

	select {
	case <-terminated:
		pe.log.Debug("process", "Process %v terminated after SIGTERM", pe.process.Cmd.Process.Pid)
		return // means process finished
	case <-time.After(500*time.Millisecond + pe.closetime):
	}

	// KILL -9
	err = pe.process.Cmd.Process.Kill()
	if err != nil {
		pe.log.Error("process", "SIGKILL unsuccessful to %v: %s", pe.process.Cmd.Process.Pid, err)
		return
	}

	select {
	case <-terminated:
		pe.log.Debug("process", "Process %v terminated after SIGKILL", pe.process.Cmd.Process.Pid)
		return // means process finished
	case <-time.After(1000 * time.Millisecond):
	}

	pe.log.Error("process", "SIGKILL did not terminate %v!", pe.process.Cmd.Process.Pid)
}

func (pe *ProcessEndpoint) Output() chan []byte {
	return pe.output
}

// 将第一个命令的输出作为第二个命令的输入
// 将第二个命令的输出作为第一个命令的输入
func (pe *ProcessEndpoint) Send(msg []byte) bool {
	pe.process.Stdin.Write(msg)
	return true
}

// 从标准输出中读取数据发给管道
func (pe *ProcessEndpoint) StartReading() {
	// 创建协程，从标准错误输出中读取数据，并打印出来
	go pe.log_stderr()
	// 从标准输出中读取
	if pe.bin {
		go pe.process_binout()
	} else {
		go pe.process_txtout()
	}
}

// 从标准输出按行读取字节数组
func (pe *ProcessEndpoint) process_txtout() {
	bufin := bufio.NewReader(pe.process.Stdout)
	for {
		buf, err := bufin.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				pe.log.Error("process", "Unexpected error while reading STDOUT from process: %s", err)
			} else {
				pe.log.Debug("process", "Process STDOUT closed")
			}
			break
		}
		pe.output <- TrimEOL(buf)
	}
	// 命令执行结束关闭管道
	close(pe.output)
}

// 从标准输出读取二进制数据
func (pe *ProcessEndpoint) process_binout() {
	buf := make([]byte, 10*1024*1024) // 10M
	for {
		n, err := pe.process.Stdout.Read(buf)
		if err != nil {
			if err != io.EOF {
				pe.log.Error("process", "Unexpected error while reading STDOUT from process: %s", err)
			} else {
				pe.log.Debug("process", "Process STDOUT closed")
			}
			break
		}
		// 复制切片到阻塞管道
		pe.output <- append(make([]byte, 0, n), buf[:n]...) // cloned buffer
	}
	// 命令执行结束关闭管道
	close(pe.output)
}

// 打印标准错误
func (pe *ProcessEndpoint) log_stderr() {
	bufStderr := bufio.NewReader(pe.process.Stderr)
	for {
		// 返回在第一次出现传入字节前的字节，没有数据则阻塞
		buf, err := bufStderr.ReadSlice('\n')
		if err != nil {
			if err != io.EOF {
				pe.log.Error("process", "Unexpected error while reading STDERR from process: %s", err)
			} else {
				pe.log.Debug("process", "Process STDERR closed")
			}
			break
		}
		pe.log.Error("stderr", "%s", string(TrimEOL(buf)))
	}
}
