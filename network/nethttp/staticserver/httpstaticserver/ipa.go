package httpstaticserver

import (
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/ios"
	"github.com/gorilla/mux"
)

func (s *HTTPStaticServer) hPlist(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	// rename *.plist to *.ipa
	if filepath.Ext(path) == ".plist" {
		path = path[0:len(path)-6] + ".ipa"
	}

	// 读取压缩包中的 Info.plist 文件
	relPath := s.getRealPath(r)
	plinfo, err := ios.ParseIPA(relPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := &url.URL{
		Scheme: scheme,
		Host:   r.Host,
	}
	data, err := ios.GenerateDownloadPlist(baseURL, path, plinfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/xml")
	w.Write(data)
}
