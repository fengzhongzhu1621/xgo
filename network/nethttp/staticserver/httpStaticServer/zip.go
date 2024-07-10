package httpstaticserver

import (
	"mime"
	"net/http"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/crypto/zipfile"
	"github.com/gorilla/mux"
)

func (s *HTTPStaticServer) hZip(w http.ResponseWriter, r *http.Request) {
	// 获得静态路由执行的文件所在的绝对路径
	file_path := s.getRealPath(r)
	// 压缩文件返回响应
	zipfile.CompressToZip(w, file_path)
}

func (s *HTTPStaticServer) hUnzip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipPath, path := vars["zip_path"], vars["path"]
	// 获取文件路径中文件的扩展名(包含 .)
	// 根据文件扩展名获取相应的MIME类型
	ctype := mime.TypeByExtension(filepath.Ext(path))
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}
	// 从压缩包中解压指定路径的文件
	err := zipfile.ExtractFromZip(filepath.Join(s.Root, zipPath), path, w)
	if err != nil {
		// 返回 500 的响应
		http.Error(w, err.Error(), 500)
		return
	}
}
