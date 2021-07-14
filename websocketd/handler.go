package websocketd

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
	"time"
	. "xgo/utils/net_utils"
)


var ScriptNotFoundError = errors.New("script not found")

// WebsocketdHandler is a single request information and processing structure, it handles WS requests out of all that daemon can handle (static, cgi, devconsole)
type WebsocketdHandler struct {
	server *WebsocketdServer

	Id string
	*RemoteInfo
	*URLInfo // TODO: I cannot find where it's used except in one single place as URLInfo.FilePath
	Env      []string

	command string
}


// NewWebsocketdHandler constructs the struct and parses all required things in it...
func NewWebsocketdHandler(s *WebsocketdServer, req *http.Request, log *LogScope) (wsh *WebsocketdHandler, err error) {
	wsh = &WebsocketdHandler{server: s, Id: generateId()}
	log.Associate("id", wsh.Id)

	wsh.RemoteInfo, err = GetRemoteInfo(req.RemoteAddr, s.Config.ReverseLookup)
	if err != nil {
		log.Error("session", "Could not understand remote address '%s': %s", req.RemoteAddr, err)
		return nil, err
	}
	log.Associate("remote", wsh.RemoteInfo.Host)

	wsh.URLInfo, err = GetURLInfo(req.URL.Path, s.Config)
	if err != nil {
		log.Access("session", "NOT FOUND: %s", err)
		return nil, err
	}

	wsh.command = s.Config.CommandName
	if s.Config.UsingScriptDir {
		wsh.command = wsh.URLInfo.FilePath
	}
	log.Associate("command", wsh.command)

	wsh.Env = createEnv(wsh, req, log)

	return wsh, nil
}

func (wsh *WebsocketdHandler) accept(ws *websocket.Conn, log *LogScope) {
	defer ws.Close()

	log.Access("session", "CONNECT")
	defer log.Access("session", "DISCONNECT")

	launched, err := launchCmd(wsh.command, wsh.server.Config.CommandArgs, wsh.Env)
	if err != nil {
		log.Error("process", "Could not launch process %s %s (%s)", wsh.command, strings.Join(wsh.server.Config.CommandArgs, " "), err)
		return
	}

	log.Associate("pid", strconv.Itoa(launched.cmd.Process.Pid))

	binary := wsh.server.Config.Binary
	process := NewProcessEndpoint(launched, binary, log)
	if cms := wsh.server.Config.CloseMs; cms != 0 {
		process.closetime += time.Duration(cms) * time.Millisecond
	}
	wsEndpoint := NewWebSocketEndpoint(ws, binary, log)

	PipeEndpoints(process, wsEndpoint)
}
