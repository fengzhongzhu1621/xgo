package httpstaticserver

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func (s *HTTPStaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func NewHTTPStaticServer(root string, noIndex bool) *HTTPStaticServer {
	// if root == "" {
	// 	root = "./"
	// }
	// root = filepath.ToSlash(root)
	root = filepath.ToSlash(filepath.Clean(root))
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}
	log.Printf("root path: %s\n", root)
	m := mux.NewRouter()
	s := &HTTPStaticServer{
		Root:  root,
		Theme: "black",
		m:     m,
		bufPool: sync.Pool{
			New: func() interface{} { return make([]byte, 32*1024) },
		},
		NoIndex: noIndex,
	}

	// 创建索引
	if !noIndex {
		go func() {
			time.Sleep(1 * time.Second)
			for {
				startTime := time.Now()
				log.Println("Started making search index")
				s.makeIndex()
				log.Printf("Completed search index in %v", time.Since(startTime))
				//time.Sleep(time.Second * 1)
				time.Sleep(time.Minute * 10)
			}
		}()
	}

	// routers for Apple *.ipa
	m.HandleFunc("/-/ipa/plist/{path:.*}", s.hPlist)
	m.HandleFunc("/-/ipa/link/{path:.*}", s.hIpaLink)

	m.HandleFunc("/{path:.*}", s.hIndex).Methods("GET", "HEAD")
	m.HandleFunc("/{path:.*}", s.hUploadOrMkdir).Methods("POST")
	m.HandleFunc("/{path:.*}", s.hDelete).Methods("DELETE")
	return s
}
