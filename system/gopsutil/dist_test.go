package gopsutil

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/shirou/gopsutil/v4/disk"
)

// 获取机器下不同磁盘分区的内容
func TestGetDiskPartitions(t *testing.T) {
	partitions, _ := disk.Partitions(false)
	// tests.PrintStruct(partitions)
	// 返回一个 []disk.PartitionStat 类型的切片，其中每个元素包含一个磁盘分区的信息，包括分区设备名、分区挂载点等。
	// [
	// 		{
	// 			"device": "/dev/disk3s1s1",
	// 			"mountpoint": "/",
	// 			"fstype": "apfs",
	// 			"opts": "ro,journaled,multilabel"
	// 		},
	// ]
	//
	for _, p := range partitions[:1] {
		// 获取该分区的使用情况
		usage, _ := disk.Usage(p.Mountpoint)
		tests.PrintStruct(usage)
		// 返回一个 *disk.UsageStat 类型的结构体，包含该磁盘分区的使用情况，如总容量、已用容量、可用容量等。
		// {
		// 			"path": "/",
		// 			"fstype": "apfs",
		// 			"total": 994662584320,
		// 			"free": 683694559232,
		// 			"used": 310968025088,
		// 			"usedPercent": 31.26366970972302,
		// 			"inodesTotal": 4290378387,
		// 			"inodesUsed": 404475,
		// 			"inodesFree": 4289973912,
		// 			"inodesUsedPercent": 0.00942749015391215
		// }
	}
}
