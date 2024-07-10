package httpStaticServer

import "path/filepath"

func (s *HTTPStaticServer) historyDirSize(dir string) int64 {
	// 从本地缓存获得目录大小
	dirInfoSize.mutex.RLock()
	size, ok := dirInfoSize.size[dir]
	dirInfoSize.mutex.RUnlock()

	if ok {
		return size
	}
	
	// 遍历索引查找索引匹配的文件大小之和
	for _, fitem := range s.indexes {
		if filepath.HasPrefix(fitem.Path, dir) {
			size += fitem.Info.Size()
		}
	}

	// 刷新路径的大小
	dirInfoSize.mutex.Lock()
	dirInfoSize.size[dir] = size
	dirInfoSize.mutex.Unlock()

	return size
}
