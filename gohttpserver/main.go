package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	accesslog "github.com/codeskyblue/go-accesslog"
	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/auth"
	"github.com/fengzhongzhu1621/xgo/network"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/fengzhongzhu1621/xgo/network/nethttp/staticserver/httpstaticserver"
	"github.com/fengzhongzhu1621/xgo/str/stringutils"
	"github.com/go-yaml/yaml"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// 解析命令行参数和配置文件
	if err := httpstaticserver.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	gcfg := httpstaticserver.GetGcfg()
	if gcfg.Debug {
		// DEBUG模式下打印配置文件的内容
		data, _ := yaml.Marshal(gcfg)
		fmt.Printf("--- config ---\n%s\n", string(data))
	}
	// 日志添加文件行号和时间
	log.SetFlags(log.Llongfile | log.LstdFlags)

	// make sure prefix matches: ^/.*[^/]$
	gcfg.Prefix = stringutils.FixPrefix(gcfg.Prefix)
	if gcfg.Prefix != "" {
		log.Printf("url prefix: %s", gcfg.Prefix)
	}

	// 创结静态服务对象, http.Handler，构造handler
	ss := httpstaticserver.NewHTTPStaticServer(gcfg.Root, gcfg.NoIndex)
	ss.Prefix = gcfg.Prefix
	ss.Theme = gcfg.Theme
	ss.Title = gcfg.Title
	ss.GoogleTrackerID = gcfg.GoogleTrackerID
	ss.Upload = gcfg.Upload
	ss.Delete = gcfg.Delete
	ss.AuthType = gcfg.Auth.Type
	ss.DeepPathMaxDepth = gcfg.DeepPathMaxDepth

	if gcfg.PlistProxy != "" {
		u, err := url.Parse(gcfg.PlistProxy)
		if err != nil {
			log.Fatal(err)
		}
		u.Scheme = "https"
		ss.PlistProxy = u.String()
	}
	if ss.PlistProxy != "" {
		log.Printf("plistproxy: %s", strconv.Quote(ss.PlistProxy))
	}
	var hdlr http.Handler = ss

	// 设置日志处理
	var logger = httpstaticserver.GetLogger()
	hdlr = accesslog.NewLoggingHandler(hdlr, logger)

	// 添加认证处理器
	// HTTP Basic Authentication
	userpass := strings.SplitN(gcfg.Auth.HTTP, ":", 2)
	switch gcfg.Auth.Type {
	case "http":
		if len(userpass) == 2 {
			user, pass := userpass[0], userpass[1]
			hdlr = httpauth.SimpleBasicAuth(user, pass)(hdlr)
		}
	case "openid":
		auth.HandleOpenID(gcfg.Auth.OpenID, false) // FIXME(ssx): set secure default to false
		// case "github":
		// 	handleOAuth2ID(gcfg.Auth.Type, gcfg.Auth.ID, gcfg.Auth.Secret) // FIXME(ssx): set secure default to false
	case "oauth2-proxy":
		auth.HandleOauth2()
	}

	// 添加 CORS headers
	hdlr = nethttp.Cors(hdlr)

	// 添加代理头
	if gcfg.XHeaders {
		hdlr = handlers.ProxyHeaders(hdlr)
	}

	// 添加路由
	mainRouter := mux.NewRouter()
	router := mainRouter
	if gcfg.Prefix != "" {
		router = mainRouter.PathPrefix(gcfg.Prefix).Subrouter()
		mainRouter.Handle(gcfg.Prefix, hdlr)
		mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, gcfg.Prefix, http.StatusTemporaryRedirect)
		})
	}

	// 添加静态资源路由， 从请求的 URL 路径中删除指定的前缀
	// 将匹配的路径 /-/assets/ -> assets/ 交给 http.FileServer 处理
	// 例子 http://localhost:8000/-/assets/css/style.css
	router.PathPrefix("/-/assets/").Handler(http.StripPrefix(gcfg.Prefix+"/-/", http.FileServer(xgo.Assets)))
	// 添加版本路由
	// 例子 http://localhost:8000/-/sysinfo
	var version = httpstaticserver.GetVersion()
	router.HandleFunc("/-/sysinfo", func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(map[string]interface{}{
			"version": version,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	})
	// 添加业务逻辑路由
	// http://localhost:8000/assets/css
	router.PathPrefix("/").Handler(hdlr)

	// 格式化服务器地址
	if gcfg.Addr == "" {
		gcfg.Addr = fmt.Sprintf(":%d", gcfg.Port)
	}
	if !strings.Contains(gcfg.Addr, ":") {
		gcfg.Addr = ":" + gcfg.Addr
	}
	_, port, _ := net.SplitHostPort(gcfg.Addr)
	log.Printf("listening on %s, local address http://%s:%s\n", strconv.Quote(gcfg.Addr), network.GetLocalIP(), port)

	// 启动服务
	srv := &http.Server{
		Handler: mainRouter,
		Addr:    gcfg.Addr,
	}
	var err error
	if gcfg.Key != "" && gcfg.Cert != "" {
		err = srv.ListenAndServeTLS(gcfg.Cert, gcfg.Key)
	} else {
		err = srv.ListenAndServe()
	}
	log.Fatal(err)
}
