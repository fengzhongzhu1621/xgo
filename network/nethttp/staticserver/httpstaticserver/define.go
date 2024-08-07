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
	Conf            *os.File `yaml:"-"`    // 配置文件路径
	Addr            string   `yaml:"addr"` // 静态服务的地址
	Port            int      `yaml:"port"` // 静态服务端口号
	Root            string   `yaml:"root"` // 静态资源根目录
	Prefix          string   `yaml:"prefix"`
	HTTPAuth        string   `yaml:"httpauth"`
	Cert            string   `yaml:"cert"`
	Key             string   `yaml:"key"`
	Theme           string   `yaml:"theme"` // 样式主题
	XHeaders        bool     `yaml:"xheaders"`
	Upload          bool     `yaml:"upload"` // 是否有文件上传权限
	Delete          bool     `yaml:"delete"` // 是否有文件删除权限
	PlistProxy      string   `yaml:"plistproxy"`
	Title           string   `yaml:"title"` // 静态服务的标题
	Debug           bool     `yaml:"debug"`
	GoogleTrackerID string   `yaml:"google-tracker-id"`
	Auth            struct {
		Type   string `yaml:"type"`   // openid|http|github
		OpenID string `yaml:"openid"` // openid 认证地址
		HTTP   string `yaml:"http"`
		ID     string `yaml:"id"`     // for oauth2
		Secret string `yaml:"secret"` // for oauth2
	} `yaml:"auth"`
	DeepPathMaxDepth int  `yaml:"deep-path-max-depth"` // 路径搜索深度
	NoIndex          bool `yaml:"no-index"`            // 进程启动时是否创建索引
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

func GetGcfg() Configure {
	return gcfg
}

func GetLogger() httpLogger {
	return logger
}

func GetVersion() string {
	return VERSION
}

func GetYamlConf() string {
	return YAMLCONF
}
