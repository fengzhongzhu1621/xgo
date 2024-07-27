package httpstaticserver

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fengzhongzhu1621/xgo/file"
	"github.com/fengzhongzhu1621/xgo/tpl"
	"github.com/gorilla/mux"
)

// 生成所有静态文件的索引配置，索引由相对路径和文件的元数据构造
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

// findIndex 根据相对路径查找查找索引
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

// hIndex 显示索引页面
func (s *HTTPStaticServer) hIndex(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	// 获得静态路由执行的文件所在的绝对路径
	realPath := s.getRealPath(r)
	log.Printf("GET hIndex path = %s realPath = %s raw = %s json = %s", path, realPath, r.FormValue("raw"), r.FormValue("json"))

	if r.FormValue("json") == "true" {
		s.hJSONList(w, r)
		return
	}

	if r.FormValue("op") == "info" {
		s.hInfo(w, r)
		return
	}

	if r.FormValue("op") == "archive" {
		s.hZip(w, r)
		return
	}

	if r.FormValue("raw") == "false" || file.IsDir(realPath) {
		// 首页访问
		if r.Method == "HEAD" {
			return
		}
		tpl.RenderHTML(w, "assets/index.html", s)
	} else {
		if filepath.Base(path) == YAMLCONF {
			auth := s.readAccessConf(realPath)
			if !auth.Delete {
				http.Error(w, "Security warning, not allowed to read", http.StatusForbidden)
				return
			}
		}
		if r.FormValue("download") == "true" {
			w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(path)))
		}
		http.ServeFile(w, r, realPath)
	}
}
