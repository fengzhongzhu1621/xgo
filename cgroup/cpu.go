// 获取 Docker 容器（或 Linux cgroup）的 CPU 使用率及相关 CPU 配置信息
package cgroup

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/file"
)

const (
	// 每个 CPU 核的使用时间（按核分开记录，单位纳秒？通常是累加值）。
	cpuFile = "/sys/fs/cgroup/cpuacct/cpuacct.usage_percpu"
	// 用于计算容器被限制的 CPU 时间配额（cfs_quota_us / cfs_period_us）。
	quotaFile = "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
	// 用于计算容器被限制的 CPU 时间配额（cfs_quota_us / cfs_period_us）。
	periodFile = "/sys/fs/cgroup/cpu/cpu.cfs_period_us"
	// 容器可以使用的 CPU 核编号，支持格式如 "0,8-12,60-63"。
	cpuSetFile = "/sys/fs/cgroup/cpuset/cpuset.cpus"
	// 容器总的 CPU 使用时间（所有核累计，单位纳秒）。
	cpuActUsageFile = "/sys/fs/cgroup/cpuacct/cpuacct.usage"
)

var (
	// 通过 getconf CLK_TCK 获取，通常是 100，表示每秒时钟滴答数（用于某些时间转换，但当前代码未直接使用）。
	cpuTick, _ = getCPUTick()
	// cores 宿主机总核数
	cores, _ = GetCoreCount()
	// 容器被 cgroup 限制可用的 CPU 配额（计算方式为 quota/period，代表每秒可用的 CPU 时间比例）
	// 容器分配可以使用的核数
	limitedCores, _ = GetLimitedCoreCount()
	// 自定义错误，表示获取 CPU 核数出错
	errCores = errors.New("Error CPU Cores")
)

// GetDockerCPUUsage 获取 interval 时间间隔内容器 cpu 的利用率
// 计算容器在某个时间间隔（比如 1s）内的 CPU 使用率。
func GetDockerCPUUsage(interval time.Duration) (usage float64, err error) {
	if interval <= 0 {
		return
	}

	// 获取 起始时刻的容器总 CPU 使用时间
	preCPUTotal, err := GetContainerCPUTotal()
	if err != nil {
		return usage, err
	}

	// 这里阻塞 interval 时间
	time.Sleep(interval)

	// 获取 结束时刻的容器总 CPU 使用时间
	postCPUTotal, err := GetContainerCPUTotal()
	if err != nil {
		return usage, err
	}

	// 计算两次采样之间的 CPU 时间差 usedCPU
	usedCPU := float64(postCPUTotal - preCPUTotal)

	// 容器 cpu 配额比
	if cores == 0 {
		return usage, errCores
	}

	// CPU 时间差 / (时间间隔 * 每秒可用的 CPU 时间配额)，得到的是一个 利用率比例（如 0.5 表示 50%）
	usage = usedCPU / (float64(interval) * limitedCores)

	return
}

// GetContainerCPUTotal 获取容器 cpu 使用时间
func GetContainerCPUTotal() (usage uint64, err error) {
	return file.ReadUint64FromFile(cpuActUsageFile)
}

// GetCoreCount 获取宿主机总共可用的 CPU 核数
func GetCoreCount() (num uint64, err error) {
	dat, err := file.ReadFromFile(cpuFile)
	if err != nil {
		return 0, err
	}

	items := strings.Split(dat, " ")
	return uint64(len(items)), nil
}

// GetLimitedCoreCount 获取容器 cpu 配额
func GetLimitedCoreCount() (mum float64, err error) {
	// 读取每个 period 时间间隔内可以使用的 cpu 时间
	quota, err := file.ReadInt64FromFile(quotaFile)
	if err != nil {
		return 0.0, err
	}

	// quota 为 -1 表示无限制，直接取可用的 cpu 核数
	if quota == -1 {
		return getValidCPUSet()
	}

	period, err := file.ReadInt64FromFile(periodFile)
	if err != nil {
		return 0.0, err
	}
	if period <= 0 {
		return 0.0, errors.New("invalid period num")
	}

	return float64(quota) / float64(period), nil
}

// getValidCPUSet 获取容器可以使用的 cpu 核
// 文件的内容存储格式例如：0,8-12,60-63
func getValidCPUSet() (ret float64, err error) {
	dat, err := file.ReadFromFile(cpuSetFile)
	if err != nil {
		return 0, err
	}

	validCores := 0
	items := strings.Split(dat, ",")
	for _, item := range items {
		cores := strings.Split(item, "-")

		switch len(cores) {
		case 1:
			// 直接指定某个核可用
			validCores++
		case 2:
			// 计算区间核数
			start, _ := strconv.Atoi(cores[0])
			end, _ := strconv.Atoi(cores[1])
			validCores += end - start + 1
		default:
			return 0.0, errors.New("Invalid cpu set formmat")
		}

	}

	return float64(validCores), nil
}

// getCPUTick 获取 CPU 调度周期（cpu 时钟）, 一般默认 100ms
func getCPUTick() (cnt int, err error) {
	out, err := exec.Command("getconf", "CLK_TCK").Output()
	if err != nil {
		return cnt, err
	}

	cnt, err = strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return cnt, err
	}

	return
}
