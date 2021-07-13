package websocketd

import (
	"fmt"
	"net/http"
	"xgo/utils/string_utils"
	. "xgo/env"
	. "xgo/log"
	. "xgo/utils/net_utils"
)

const (
	gatewayInterface = "websocketd-CGI/0.1"
)

func createEnv(handler *WebsocketdHandler, req *http.Request, log *LogScope) []string {
	headers := req.Header

	url := req.URL

	// 获得主机和端口
	serverName, serverPort, err :=TellHostPort(req.Host, handler.server.Config.Ssl)
	if err != nil {
		// This does mean that we cannot detect port from Host: header... Just keep going with "", guessing is bad.
		log.Debug("env", "Host port detection error: %s", err)
		serverPort = ""
	}

	standardEnvCount := 20
	if handler.server.Config.Ssl {
		standardEnvCount += 1
	}

	parentLen := len(handler.server.Config.ParentEnv)
	env := make([]string, 0, len(headers)+standardEnvCount+parentLen+len(handler.server.Config.Env))

	// This variable could be rewritten from outside
	env = AppendEnv(env, "SERVER_SOFTWARE", handler.server.Config.ServerSoftware)

	parentStarts := len(env)
	env = append(env, handler.server.Config.ParentEnv...)

	// IMPORTANT ---> Adding a header? Make sure standardEnvCount (above) is up to date.

	// Standard CGI specification headers.
	// As defined in http://tools.ietf.org/html/rfc3875
	env = AppendEnv(env, "REMOTE_ADDR", handler.RemoteInfo.Addr)
	env = AppendEnv(env, "REMOTE_HOST", handler.RemoteInfo.Host)
	env = AppendEnv(env, "SERVER_NAME", serverName)	// IP
	env = AppendEnv(env, "SERVER_PORT", serverPort)	// port
	env = AppendEnv(env, "SERVER_PROTOCOL", req.Proto)
	env = AppendEnv(env, "GATEWAY_INTERFACE", gatewayInterface)
	env = AppendEnv(env, "REQUEST_METHOD", req.Method)
	env = AppendEnv(env, "SCRIPT_NAME", handler.URLInfo.ScriptPath)
	env = AppendEnv(env, "PATH_INFO", handler.URLInfo.PathInfo)
	env = AppendEnv(env, "PATH_TRANSLATED", url.Path)
	env = AppendEnv(env, "QUERY_STRING", url.RawQuery)

	// Not supported, but we explicitly clear them so we don't get leaks from parent environment.
	env = AppendEnv(env, "AUTH_TYPE", "")
	env = AppendEnv(env, "CONTENT_LENGTH", "")
	env = AppendEnv(env, "CONTENT_TYPE", "")
	env = AppendEnv(env, "REMOTE_IDENT", "")
	env = AppendEnv(env, "REMOTE_USER", "")

	// Non standard, but commonly used headers.
	env = AppendEnv(env, "UNIQUE_ID", handler.Id) // Based on Apache mod_unique_id.
	env = AppendEnv(env, "REMOTE_PORT", handler.RemoteInfo.Port)
	env = AppendEnv(env, "REQUEST_URI", url.RequestURI()) // e.g. /foo/blah?a=b

	// The following variables are part of the CGI specification, but are optional
	// and not set by websocketd:
	//
	//   AUTH_TYPE, REMOTE_USER, REMOTE_IDENT
	//     -- Authentication left to the underlying programs.
	//
	//   CONTENT_LENGTH, CONTENT_TYPE
	//     -- makes no sense for WebSocket connections.
	//
	//   SSL_*
	//     -- SSL variables are not supported, HTTPS=on added for websocketd running with --ssl

	if handler.server.Config.Ssl {
		env = AppendEnv(env, "HTTPS", "on")
	}

	if log.MinLevel == LogDebug {
		for i, v := range env {
			if i >= parentStarts && i < parentLen+parentStarts {
				log.Debug("env", "Parent envvar: %v", v)
			} else {
				log.Debug("env", "Std. variable: %v", v)
			}
		}
	}

	for k, hdrs := range headers {
		header := fmt.Sprintf("HTTP_%s", string_utils.HeaderDashToUnderscore.Replace(k))
		env = AppendEnv(env, header, hdrs...)
		log.Info("env", "Header variable %s", env[len(env)-1])
	}

	for _, v := range handler.server.Config.Env {
		env = append(env, v)
		log.Debug("env", "External variable: %s", v)
	}

	return env
}
