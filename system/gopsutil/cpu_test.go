package gopsutil

import (
	"fmt"
	"log"
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

func TestCpuPercent(t *testing.T) {
	pid := 9870
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		log.Fatalf("无法获取进程 %d 的信息: %v", pid, err)
	}

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		log.Printf("获取 CPU 使用率失败: %v", err)
		return
	}

	if cpuPercent != 0.0 {
		fmt.Printf("CPU 使用率: %.2f%%\n", cpuPercent) // CPU 使用率: 1.62%
	}
}
