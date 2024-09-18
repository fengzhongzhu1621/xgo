//go:build linux
// +build linux

package taskset

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"unsafe"
)

// 使用 linux 下的 taskset 命令查询和设置进程的 cpu 亲和性
func TestTaskSet(t *testing.T) {
	pid := os.Getpid()

	// 获得进程 CPU 亲和性掩码
	cmd := exec.Command("taskset", "-p", fmt.Sprintf("%d", pid))
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// 运行时设置 CPU 亲和性掩码
	cmd = exec.Command("taskset", "-p", "0,1", fmt.Sprintf("%d", pid))
	out, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", out)
}

// 使用系统调用设置 CPU 亲和性
func TestSchedSetAffinity(t *testing.T) {
	var mask uintptr

	// 获取当前进程的 CPU 亲和性
	if _, _, err := syscall.RawSyscall(syscall.SYS_SCHED_GETAFFINITY,
		0,
		uintptr(unsafe.Sizeof(mask)),
		uintptr(unsafe.Pointer(&mask)),
	); err != 0 {
		fmt.Println("failed to get CPU affinity: ", err)
		return
	}
	fmt.Println("current CPU affinity: ", mask)

	// 设置当前进程的 CPU 亲和性为 CPU0 和 CPU1
	mask = 3
	if _, _, err := syscall.RawSyscall(syscall.SYS_SCHED_SETAFFINITY,
		0,
		uintptr(unsafe.Sizeof(mask)),
		uintptr(unsafe.Pointer(&mask)),
	); err != 0 {
		fmt.Println("failed to set CPU affinity: ", err)
		return
	}
	fmt.Println("new CPU affinity: ", mask)
}
