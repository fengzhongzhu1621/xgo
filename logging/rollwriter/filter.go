package rollwriter

import (
	"strings"
	"time"
)

// filterByMaxBackups 根据最大备份数过滤冗余文件
// 参数：files - 文件列表，remove - 待删除文件列表指针，maxBackups - 最大备份数
// 返回值：保留的文件列表
func filterByMaxBackups(files []logInfo, remove *[]logInfo, maxBackups int) []logInfo {
	if maxBackups == 0 || len(files) < maxBackups { // 如果不需要限制或文件数未超过限制
		return files
	}
	var remaining []logInfo
	preserved := make(map[string]bool) // 用于跟踪已保留的文件名（去除压缩后缀）

	for _, f := range files {
		// 去除压缩后缀的文件名，用于判断是否已保留过该文件
		fn := strings.TrimSuffix(f.Name(), compressSuffix)
		preserved[fn] = true

		// 如果已保留的文件数超过限制
		if len(preserved) > maxBackups {
			// 添加到待删除列表
			*remove = append(*remove, f)
		} else {
			// 添加到保留列表
			remaining = append(remaining, f)
		}
	}

	// 返回保留的文件信息列表
	return remaining
}

// filterByMaxAge 根据最大保存天数过滤过期文件
// 参数：files - 文件列表，remove - 待删除文件列表指针，maxAge - 最大保存天数
// 返回值：保留的文件列表
func filterByMaxAge(files []logInfo, remove *[]logInfo, maxAge int) []logInfo {
	if maxAge <= 0 { // 如果不需要按时间过滤
		return files
	}

	var remaining []logInfo
	diff := time.Duration(int64(24*time.Hour) * int64(maxAge)) // 计算总的时间间隔
	cutoff := time.Now().Add(-1 * diff)                        // 计算截止时间点
	for _, f := range files {
		if f.timestamp.Before(cutoff) { // 如果文件时间早于截止时间
			// 添加到待删除列表
			*remove = append(*remove, f)
		} else {
			// 添加到待删除列表
			remaining = append(remaining, f)
		}
	}
	return remaining
}

// filterByCompressExt 根据压缩扩展名过滤需要压缩的文件
// 参数：files - 文件列表，compress - 待压缩文件列表指针，needCompress - 是否需要压缩
func filterByCompressExt(files []logInfo, compress *[]logInfo, needCompress bool) {
	if !needCompress { // 如果不需要压缩
		return
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), compressSuffix) { // 如果文件没有压缩后缀
			*compress = append(*compress, f) // 添加到待压缩列表
		}
	}
}
