package websocketd

import (
	"github.com/gorilla/websocket"
	. "xgo/utils/net_utils"
)


// WebsocketdHandler is a single request information and processing structure, it handles WS requests out of all that daemon can handle (static, cgi, devconsole)
type WebsocketdHandler struct {
	server *WebsocketdServer

	Id string
	*RemoteInfo
	*URLInfo // TODO: I cannot find where it's used except in one single place as URLInfo.FilePath
	Env      []string

	command string
}

