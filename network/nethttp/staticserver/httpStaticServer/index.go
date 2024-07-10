package httpStaticServer

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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


// 查找索引，根据相对路径查找
func (s *HTTPStaticServer) findIndex(text string) []IndexFileItem {
	ret := make([]IndexFileItem, 0)
	// 遍历索引配置
	for _, item := range s.indexes {
		ok := true
		// search algorithm, space for AND
		// 将一个字符串按照空白字符（如空格、制表符、换行符等）进行分割，并返回一个包含分割后的子字符串的切片（slice）
		for _, keyword := range strings.Fields(text) {
			needContains := true
			// 去掉开头的 -，排除类似于 exclude
			if strings.HasPrefix(keyword, "-") {
				needContains = false
				keyword = keyword[1:]
			}
			if keyword == "" {
				continue
			}
			// 检查一个字符串是否包含另一个子字符串。
			ok = (needContains == strings.Contains(strings.ToLower(item.Path), strings.ToLower(keyword)))
			if !ok {
				break
			}
		}
		if ok {
			ret = append(ret, item)
		}
	}
	return ret
}