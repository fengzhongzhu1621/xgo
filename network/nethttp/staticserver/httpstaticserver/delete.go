package httpstaticserver

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func (s *HTTPStaticServer) hDelete(w http.ResponseWriter, req *http.Request) {
	path := mux.Vars(req)["path"]
	// 获得静态路由执行的文件所在的绝对路径
	realPath := s.getRealPath(req)
	// 从路径所在的目录和祖先目录读取静态文件关联的配置文件，转换为结构体
	auth := s.readAccessConf(realPath)
	if !auth.CanDelete(req) {
		http.Error(w, "Delete forbidden", http.StatusForbidden)
		return
	}

	// TODO: path safe check
	err := os.RemoveAll(realPath)
	if err != nil {
		pathErr, ok := err.(*os.PathError)
		if ok {
			http.Error(w, pathErr.Op+" "+path+": "+pathErr.Err.Error(), 500)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	w.Write([]byte("Success"))
}
