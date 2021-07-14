package websocketd

import (
	"errors"
	"os"
	. "xgo/log"
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
