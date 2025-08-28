package gopsutil

import (
	"github.com/shirou/gopsutil/load"
)

// GetLoadAvg 获取 load 负载，包含 load1，load5，load15
func GetLoadAvg() (*load.AvgStat, error) {
	loadavg, err := load.Avg()
	if err != nil {
		return &load.AvgStat{}, err
	}
	return loadavg, nil
}
