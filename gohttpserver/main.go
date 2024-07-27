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
	if err := httpstaticserver.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	gcfg := httpstaticserver.GetGcfg()
	if gcfg.Debug {
		data, _ := yaml.Marshal(gcfg)
		fmt.Printf("--- config ---\n%s\n", string(data))
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// make sure prefix matches: ^/.*[^/]$
	gcfg.Prefix = stringutils.FixPrefix(gcfg.Prefix)
	if gcfg.Prefix != "" {
		log.Printf("url prefix: %s", gcfg.Prefix)
	}

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
	var logger = httpstaticserver.GetLogger()
	hdlr = accesslog.NewLoggingHandler(hdlr, logger)

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

	// CORS
	hdlr = nethttp.Cors(hdlr)

	if gcfg.XHeaders {
		hdlr = handlers.ProxyHeaders(hdlr)
	}

	mainRouter := mux.NewRouter()
	router := mainRouter
	if gcfg.Prefix != "" {
		router = mainRouter.PathPrefix(gcfg.Prefix).Subrouter()
		mainRouter.Handle(gcfg.Prefix, hdlr)
		mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, gcfg.Prefix, http.StatusTemporaryRedirect)
		})
	}

	var version = httpstaticserver.GetVersion()
	router.PathPrefix("/-/assets/").Handler(http.StripPrefix(gcfg.Prefix+"/-/", http.FileServer(xgo.Assets)))
	router.HandleFunc("/-/sysinfo", func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(map[string]interface{}{
			"version": version,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	})
	router.PathPrefix("/").Handler(hdlr)

	if gcfg.Addr == "" {
		gcfg.Addr = fmt.Sprintf(":%d", gcfg.Port)
	}
	if !strings.Contains(gcfg.Addr, ":") {
		gcfg.Addr = ":" + gcfg.Addr
	}
	_, port, _ := net.SplitHostPort(gcfg.Addr)
	log.Printf("listening on %s, local address http://%s:%s\n", strconv.Quote(gcfg.Addr), network.GetLocalIP(), port)

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
