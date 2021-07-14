package websocketd

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/cgi"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	. "xgo/log"
	. "xgo/utils/net_utils"
)

var ForkNotAllowedError = errors.New("too many forks active")

// WebsocketdServer presents http.Handler interface for requests libwebsocketd is handling.
type WebsocketdServer struct {
	Config *Config
	Log    *LogScope
	forks  chan byte
}

// NewWebsocketdServer creates WebsocketdServer struct with pre-determined config, logscope and maxforks limit
func NewWebsocketdServer(config *Config, log *LogScope, maxforks int) *WebsocketdServer {
	mux := &WebsocketdServer{
		Config: config,
		Log:    log,
	}
	// 设置管道的缓存大小
	if maxforks > 0 {
		mux.forks = make(chan byte, maxforks)
	}
	return mux
}

var canonicalHostname string

// 获得服务器地址
// TellURL is a helper function that changes http to https or ws to wss in case if SSL is used
func (h *WebsocketdServer) TellURL(scheme, host, path string) string {
	if len(host) > 0 && host[0] == ':' {
		if canonicalHostname == "" {
			var err error
			// 返回内核提供的主机名
			canonicalHostname, err = os.Hostname()
			if err != nil {
				canonicalHostname = "UNKNOWN"
			}
		}
		host = canonicalHostname + host
	}
	if h.Config.Ssl {
		return scheme + "s://" + host + path
	}
	return scheme + "://" + host + path
}

// 通过管道创建一个fock标识
func (h *WebsocketdServer) noteForkCreated() error {
	// note that forks can be nil since the construct could've been created by
	// someone who is not using NewWebsocketdServer
	if h.forks != nil {
		select {
		case h.forks <- 1:
			// 通过缓存区的方式限制focks的数量
			return nil
		default:
			// 管道关闭时抛出异常
			return ForkNotAllowedError
		}
	} else {
		return nil
	}
}

// fock完成时从管道中删除
func (h *WebsocketdServer) noteForkCompled() {
	if h.forks != nil { // see comment in noteForkCreated
		select {
		case <-h.forks:
			return
		default:
			// This could only happen if the completion handler called more times than creation handler above
			// Code should be audited to not allow this to happen, it's desired to have test that would
			// make sure this is impossible but it is not exist yet.
			panic("Cannot deplet number of allowed forks, something is not right in code!")
		}
	}
}

// ServeHTTP muxes between WebSocket handler, CGI handler, DevConsole, Static HTML or 404.
func (h *WebsocketdServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := h.Log.NewLevel(h.Log.LogFunc)
	log.Associate("url", h.TellURL("http", req.Host, req.RequestURI))

	// 存在命令和可执行脚本
	if h.Config.CommandName != "" || h.Config.UsingScriptDir {
		hdrs := req.Header
		upgradeRe := regexp.MustCompile(`(?i)(^|[,\s])Upgrade($|[,\s])`)
		// WebSocket, limited to size of h.forks
		if strings.ToLower(hdrs.Get("Upgrade")) == "websocket" && upgradeRe.MatchString(hdrs.Get("Connection")) {
			if h.noteForkCreated() == nil {
				defer h.noteForkCompled()

				// start figuring out if we even need to upgrade
				handler, err := NewWebsocketdHandler(h, req, log)
				if err != nil {
					if err == ScriptNotFoundError {
						log.Access("session", "NOT FOUND: %s", err)
						http.Error(w, "404 Not Found", 404)
					} else {
						log.Access("session", "INTERNAL ERROR: %s", err)
						http.Error(w, "500 Internal Server Error", 500)
					}
					return
				}

				var headers http.Header
				if len(h.Config.Headers)+len(h.Config.HeadersWs) > 0 {
					headers = http.Header(make(map[string][]string))
					PushHeaders(headers, h.Config.Headers)
					PushHeaders(headers, h.Config.HeadersWs)
				}

				upgrader := &websocket.Upgrader{
					HandshakeTimeout: h.Config.HandshakeTimeout,
					CheckOrigin: func(r *http.Request) bool {
						// backporting previous checkorigin for use in gorilla/websocket for now
						err := checkOrigin(req, h.Config, log)
						return err == nil
					},
				}
				conn, err := upgrader.Upgrade(w, req, headers)
				if err != nil {
					log.Access("session", "Unable to Upgrade: %s", err)
					http.Error(w, "500 Internal Error", 500)
					return
				}

				// old func was used in x/net/websocket style, we reuse it here for gorilla/websocket
				handler.accept(conn, log)
				return

			} else {
				log.Error("http", "Max of possible forks already active, upgrade rejected")
				http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			}
			return
		}
	}

	PushHeaders(w.Header(), h.Config.HeadersHTTP)

	// Dev console (if enabled)
	if h.Config.DevConsole {
		log.Access("http", "DEVCONSOLE")
		content := ConsoleContent
		content = strings.Replace(content, "{{license}}", License, -1)
		content = strings.Replace(content, "{{addr}}", h.TellURL("ws", req.Host, req.RequestURI), -1)
		http.ServeContent(w, req, ".html", h.Config.StartupTime, strings.NewReader(content))
		return
	}

	// CGI scripts, limited to size of h.forks
	if h.Config.CgiDir != "" {
		filePath := path.Join(h.Config.CgiDir, fmt.Sprintf(".%s", filepath.FromSlash(req.URL.Path)))
		if fi, err := os.Stat(filePath); err == nil && !fi.IsDir() {

			log.Associate("cgiscript", filePath)
			if h.noteForkCreated() == nil {
				defer h.noteForkCompled()

				// Make variables to supplement cgi... Environ it uses will show empty list.
				envlen := len(h.Config.ParentEnv)
				cgienv := make([]string, envlen+1)
				if envlen > 0 {
					copy(cgienv, h.Config.ParentEnv)
				}
				cgienv[envlen] = "SERVER_SOFTWARE=" + h.Config.ServerSoftware
				cgiHandler := &cgi.Handler{
					Path: filePath,
					Env: []string{
						"SERVER_SOFTWARE=" + h.Config.ServerSoftware,
					},
				}
				log.Access("http", "CGI")
				cgiHandler.ServeHTTP(w, req)
			} else {
				log.Error("http", "Fork not allowed since maxforks amount has been reached. CGI was not run.")
				http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			}
			return
		}
	}

	// Static files
	if h.Config.StaticDir != "" {
		handler := http.FileServer(http.Dir(h.Config.StaticDir))
		log.Access("http", "STATIC")
		handler.ServeHTTP(w, req)
		return
	}

	// 404
	log.Access("http", "NOT FOUND")
	http.NotFound(w, req)
}
