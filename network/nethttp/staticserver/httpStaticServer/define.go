package httpStaticServer

import (
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// 索引配置
type IndexFileItem struct {
	Path string // 静态文件相对于 Root 的路径
	Info os.FileInfo // 静态文件的元数据信息
}

type Directory struct {
	size  map[string]int64
	mutex *sync.RWMutex
}

type HTTPStaticServer struct {
	Root             string // 根路径
	Prefix           string // 路径前缀
	Upload           bool // 非登录用户的上传权限
	Delete           bool // 非登录用户的删除权限
	Title            string
	Theme            string
	PlistProxy       string
	GoogleTrackerID  string
	AuthType         string
	DeepPathMaxDepth int
	NoIndex          bool

	indexes []IndexFileItem // 所有静态文件的索引配置
	m       *mux.Router
	bufPool sync.Pool // use sync.Pool caching buf to reduce gc ratio
}

type FileJSONInfo struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Size    int64       `json:"size"`
	Path    string      `json:"path"`
	ModTime int64       `json:"mtime"`
	Extra   interface{} `json:"extra,omitempty"`
}

type HTTPFileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mtime"`
}

var dirInfoSize = Directory{size: make(map[string]int64), mutex: &sync.RWMutex{}}
