package httpstaticserver

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/android"
	"github.com/fengzhongzhu1621/xgo/file"
	"github.com/gorilla/mux"
)

// hInfo 返回文件或目录的元信息
func (s *HTTPStaticServer) hInfo(w http.ResponseWriter, r *http.Request) {
	// 获得静态路由执行的文件所在的绝对路径
	path := mux.Vars(r)["path"]
	relPath := s.getRealPath(r)

	log.Printf("hInfo path = %s realPath = %s", path, relPath)

	// 获得文件的元数据
	fi, err := os.Stat(relPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fji := &FileJSONInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Path:    path,
		ModTime: fi.ModTime().UnixNano() / 1e6, // 文件最新修改时间
	}
	// 获得文件的后缀
	ext := filepath.Ext(path)
	switch ext {
	case ".md":
		fji.Type = "markdown"
	case ".apk":
		fji.Type = "apk"
		fji.Extra = android.ParseApkInfo(relPath)
	case "":
		// 没有后缀默认为文件夹
		fji.Type = "dir"
	default:
		fji.Type = "text"
	}

	// 转换为 json 字符串
	data, _ := json.Marshal(fji)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// hJSONList 获得用户有权限的文件或目录元信息和当前拥有的权限配置
func (s *HTTPStaticServer) hJSONList(w http.ResponseWriter, r *http.Request) {
	// 获得静态路由执行的文件所在的绝对路径
	requestPath := mux.Vars(r)["path"]
	realPath := s.getRealPath(r)

	log.Printf("hJSONList path = %s realPath = %s", requestPath, realPath)

	// 从路径所在的目录和祖先目录读取静态文件关联的配置文件，转换为结构体
	auth := s.readAccessConf(realPath)
	// 判断登录用户是否有上传权限
	auth.Upload = auth.CanUpload(r)
	// 判断登录用户是否有删除权限
	auth.Delete = auth.CanDelete(r)
	// 获得路径搜素的最大深度
	maxDepth := s.DeepPathMaxDepth

	// path string -> info os.FileInfo
	fileInfoMap := make(map[string]os.FileInfo, 0)

	// 获取路径下所有文件和子目录的 os.FileInfo
	search := r.FormValue("search")
	if search != "" {
		// 根据相对路径查找查找索引，只获取匹配的前 50 条
		results := s.findIndex(search)
		if len(results) > 50 { // max 50
			results = results[:50]
		}
		// 获得匹配的索引
		for _, item := range results {
			if filepath.HasPrefix(item.Path, requestPath) {
				fileInfoMap[item.Path] = item.Info
			}
		}
	} else {
		// 读取指定目录下的所有条目（包括文件和子目录）。
		infos, err := os.ReadDir(realPath)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		for _, entry := range infos {
			info, err := entry.Info()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			fileInfoMap[filepath.Join(requestPath, info.Name())] = info
		}
	}

	// turn file list -> json
	lrs := make([]HTTPFileInfo, 0)
	for path, info := range fileInfoMap {
		// 忽略无权限的文件和目录的访问
		if !auth.CanAccess(info.Name()) {
			continue
		}
		lr := HTTPFileInfo{
			Name:    info.Name(),
			Path:    path,
			ModTime: info.ModTime().UnixNano() / 1e6,
		}
		if search != "" {
			name, err := filepath.Rel(requestPath, path)
			if err != nil {
				log.Println(requestPath, path, err)
			}
			lr.Name = filepath.ToSlash(name) // fix for windows
		}
		if info.IsDir() {
			// 在目录下查询子目录，获得子目录的路径
			name := file.DeepPath(realPath, info.Name(), maxDepth)
			lr.Name = name
			lr.Path = filepath.Join(filepath.Dir(path), name)
			lr.Type = "dir"
			// 获得目录的大小
			lr.Size = s.historyDirSize(lr.Path)
		} else {
			lr.Type = "file"
			// 获得文件的大小
			lr.Size = info.Size() // formatSize(info)
		}
		lrs = append(lrs, lr)
	}

	data, _ := json.Marshal(map[string]interface{}{
		"files": lrs,
		"auth":  auth,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
