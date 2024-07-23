package httpstaticserver

import (
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// 索引配置
type IndexFileItem struct {
	Path string      // 静态文件相对于 Root 的路径
	Info os.FileInfo // 静态文件的元数据信息
}

// 目录大小
type Directory struct {
	size  map[string]int64
	mutex *sync.RWMutex
}

type HTTPStaticServer struct {
	Root             string // 根路径
	Prefix           string // 路径前缀
	Upload           bool   // 非登录用户的上传权限
	Delete           bool   // 非登录用户的删除权限
	Title            string
	Theme            string
	PlistProxy       string
	GoogleTrackerID  string
	AuthType         string
	DeepPathMaxDepth int // 路径搜素的最大深度
	NoIndex          bool

	indexes []IndexFileItem // 所有静态文件的索引配置
	m       *mux.Router
	bufPool sync.Pool // use sync.Pool caching buf to reduce gc ratio
}

type FileJSONInfo struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Size    int64       `json:"size"`
	Path    string      `json:"path"`            // 文件 http请求路径
	ModTime int64       `json:"mtime"`           // 文件最新修改时间
	Extra   interface{} `json:"extra,omitempty"` // omitempty 标签选项可用于在 JSON 序列化和反序列化过程中忽略零值字段
}

type HTTPFileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mtime"`
}

var dirInfoSize = Directory{size: make(map[string]int64), mutex: &sync.RWMutex{}}

type Configure struct {
	Conf            *os.File `yaml:"-"`
	Addr            string   `yaml:"addr"`
	Port            int      `yaml:"port"`
	Root            string   `yaml:"root"`
	Prefix          string   `yaml:"prefix"`
	HTTPAuth        string   `yaml:"httpauth"`
	Cert            string   `yaml:"cert"`
	Key             string   `yaml:"key"`
	Theme           string   `yaml:"theme"`
	XHeaders        bool     `yaml:"xheaders"`
	Upload          bool     `yaml:"upload"`
	Delete          bool     `yaml:"delete"`
	PlistProxy      string   `yaml:"plistproxy"`
	Title           string   `yaml:"title"`
	Debug           bool     `yaml:"debug"`
	GoogleTrackerID string   `yaml:"google-tracker-id"`
	Auth            struct {
		Type   string `yaml:"type"` // openid|http|github
		OpenID string `yaml:"openid"`
		HTTP   string `yaml:"http"`
		ID     string `yaml:"id"`     // for oauth2
		Secret string `yaml:"secret"` // for oauth2
	} `yaml:"auth"`
	DeepPathMaxDepth int  `yaml:"deep-path-max-depth"`
	NoIndex          bool `yaml:"no-index"`
}

var (
	defaultPlistProxy = "https://plistproxy.herokuapp.com/plist"
	defaultOpenID     = "https://login.netease.com/openid"
	gcfg              = Configure{}
	logger            = httpLogger{}

	VERSION   = "unknown"
	BUILDTIME = "unknown time"
	GITCOMMIT = "unknown git commit"
	SITE      = "https://github.com/codeskyblue/gohttpserver"
)

const YAMLCONF = ".ghs.yml"
