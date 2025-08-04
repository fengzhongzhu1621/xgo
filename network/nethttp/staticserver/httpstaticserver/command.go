package httpstaticserver

import (
	"io"

	"github.com/alecthomas/kingpin"
	"gopkg.in/yaml.v2"
)

// ParseFlags 解析命令行参数和配置文件
func ParseFlags() error {
	// initial default conf
	gcfg.Root = "./"     // 静态资源根目录
	gcfg.Port = 8000     // 静态服务端口号
	gcfg.Addr = ""       // 静态服务的地址
	gcfg.Theme = "black" // 样式主题

	gcfg.PlistProxy = defaultPlistProxy
	gcfg.Auth.OpenID = defaultOpenID
	gcfg.GoogleTrackerID = "UA-81205425-2"

	gcfg.Title = "Go HTTP File Server" // 静态服务的标题
	gcfg.DeepPathMaxDepth = 5          // 路径搜索深度
	gcfg.NoIndex = false               // 进程启动时是否创建索引，默认创建索引

	// kingpin 是一个用 Go 语言编写的命令行参数解析库。它提供了一种简单、声明式的方法来定义和解析命令行参数。
	// -h
	kingpin.HelpFlag.Short('h')
	// --version
	kingpin.Version(versionMessage()) // 获得静态服务的版本信息

	// 配置文件路径
	kingpin.Flag("conf", "config file path, yaml format").FileVar(&gcfg.Conf)
	// 静态资源根目录
	kingpin.Flag("root", "root directory, default ./").Short('r').StringVar(&gcfg.Root)
	// url 前缀
	kingpin.Flag("prefix", "url prefix, eg /foo").StringVar(&gcfg.Prefix)
	// 静态服务端口号
	kingpin.Flag("port", "listen port, default 8000").IntVar(&gcfg.Port)
	// 静态服务的地址
	kingpin.Flag("addr", "listen address, eg 127.0.0.1:8000").Short('a').StringVar(&gcfg.Addr)

	// 证书路径
	kingpin.Flag("cert", "tls cert.pem path").StringVar(&gcfg.Cert)
	kingpin.Flag("key", "tls key.pem path").StringVar(&gcfg.Key)

	// 服务认证方式
	kingpin.Flag("auth-type", "Auth type <http|openid>").StringVar(&gcfg.Auth.Type)
	kingpin.Flag("auth-http", "HTTP basic auth (ex: user:pass)").StringVar(&gcfg.Auth.HTTP)
	kingpin.Flag("auth-openid", "OpenID auth identity url").StringVar(&gcfg.Auth.OpenID)

	// 样式主题
	kingpin.Flag("theme", "web theme, one of <black|green>").StringVar(&gcfg.Theme)

	// 是否有文件上传权限
	kingpin.Flag("upload", "enable upload support").BoolVar(&gcfg.Upload)
	// 是否有文件删除权限
	kingpin.Flag("delete", "enable delete support").BoolVar(&gcfg.Delete)

	// 是否添加 xheaders
	kingpin.Flag("xheaders", "used when behide nginx").BoolVar(&gcfg.XHeaders)
	// 是否开启调试
	kingpin.Flag("debug", "enable debug mode").BoolVar(&gcfg.Debug)

	// 静态服务的标题
	kingpin.Flag("title", "server title").StringVar(&gcfg.Title)
	// 路径搜索深度
	kingpin.Flag("deep-path-max-depth", "set to -1 to not combine dirs").
		IntVar(&gcfg.DeepPathMaxDepth)
	// 进程启动时是否创建索引
	kingpin.Flag("no-index", "disable indexing").BoolVar(&gcfg.NoIndex)

	// plistproxy 代理地址
	kingpin.Flag("plistproxy", "plist proxy when server is not https").
		Short('p').
		StringVar(&gcfg.PlistProxy)
	kingpin.Flag("google-tracker-id", "set to empty to disable it").StringVar(&gcfg.GoogleTrackerID)

	// 解析命令行参数
	kingpin.Parse() // first parse conf

	if gcfg.Conf != nil {
		defer func() {
			// 用于将命令行传递的参数覆盖从配置文件读取的参数
			kingpin.Parse() // command line priority high than conf
		}()
		// 读取配置文件
		ymlData, err := io.ReadAll(gcfg.Conf)
		if err != nil {
			return err
		}
		// 将配置文件转换为结构体对象
		return yaml.Unmarshal(ymlData, &gcfg)
	}
	return nil
}
