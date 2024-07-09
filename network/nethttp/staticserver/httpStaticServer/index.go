package httpStaticServer

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// 生成所有静态文件的索引配置
func (s *HTTPStaticServer) makeIndex() error {
	var indexes = make([]IndexFileItem, 0)
	// 遍历所有的静态文件
	var err = filepath.Walk(s.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("WARN: Visit path: %s error: %v", strconv.Quote(path), err)
			return filepath.SkipDir
			// return err
		}
		// 忽略目录
		if info.IsDir() {
			return nil
		}

		// 获得文件的相对路径
		path, _ = filepath.Rel(s.Root, path)
		path = filepath.ToSlash(path)
		indexes = append(indexes, IndexFileItem{path, info})
		return nil
	})

	
	s.indexes = indexes
	return err
}
