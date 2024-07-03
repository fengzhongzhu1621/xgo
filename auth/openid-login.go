package auth

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/codeskyblue/openid-go"
	"github.com/fengzhongzhu1621/xgo/network/nethttp/staticserver"
)

type M map[string]interface{}

func init() {
	// 注册用户自定义的类型，以便在编码和解码时使用
	gob.Register(&staticserver.UserInfo{})
	gob.Register(&M{})
}

// 注册 openid 登录相关的路由
func HandleOpenID(loginUrl string, secure bool) {
	// 注册登录路由
	http.HandleFunc("/-/login", func(w http.ResponseWriter, r *http.Request) {
		// 从 HTTP 请求中获取表单字段的值，值为登录成功后的跳转链接
		nextUrl := r.FormValue("next")
		// 获取 HTTP 请求头中的 Referer 字段
		referer := r.Referer()
		host := r.Host
		if nextUrl == "" && strings.Contains(referer, "://"+host) {
			nextUrl = referer
		}
		scheme := "http"
		if r.URL.Scheme != "" {
			scheme = r.URL.Scheme
		}
		// 构造登录链接
		if url, err := openid.RedirectURL(loginUrl,
			scheme+"://"+r.Host+"/-/openidcallback?next="+nextUrl, ""); err == nil {
			// 重定向到 openid 回调链接
			http.Redirect(w, r, url, 303)
		} else {
			log.Println("Should not got error here:", err)
		}
	})

	// 注册登录链接路由
	http.HandleFunc("/-/openidcallback", func(w http.ResponseWriter, r *http.Request) {
		id, err := openid.Verify("http://"+r.Host+r.URL.String(), staticserver.DiscoveryCache, staticserver.NonceStore)
		if err != nil {
			io.WriteString(w, "Authentication check failed.")
			return
		}
		// 将用户信息存放到 session 中
		session, err := staticserver.Store.Get(r, staticserver.DefaultSessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := &staticserver.UserInfo{
			Id:       id,
			Email:    r.FormValue("openid.sreg.email"),
			Name:     r.FormValue("openid.sreg.fullname"),
			NickName: r.FormValue("openid.sreg.nickname"),
		}
		session.Values["user"] = user
		if err := session.Save(r, w); err != nil {
			log.Println("session save error:", err)
		}

		nextUrl := r.FormValue("next")
		if nextUrl == "" {
			nextUrl = "/"
		}
		http.Redirect(w, r, nextUrl, 302)
	})

	// 获得登录用户信息
	http.HandleFunc("/-/user", func(w http.ResponseWriter, r *http.Request) {
		session, err := staticserver.Store.Get(r, staticserver.DefaultSessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 获得用户信息
		val := session.Values["user"]
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 将数据结构转换为 JSON 格式的字节切片（[]byte）。如果转换成功，函数会返回转换后的字节切片和一个 nil 错误；如果转换失败，则会返回一个非空的错误。
		data, _ := json.Marshal(val)
		w.Write(data)
	})

	http.HandleFunc("/-/logout", func(w http.ResponseWriter, r *http.Request) {
		// 删除 session 中的用户信息
		session, err := staticserver.Store.Get(r, staticserver.DefaultSessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		delete(session.Values, "user")
		// 用于设置 session 的过期时间，单位为秒。
		// 当 MaxAge 设置为正整数时，表示 session 在该时间后会被自动删除。
		// 当 MaxAge 设置为负数时，表示 session 是一个临时 session，不会被持久化到存储引擎中，而是在每次请求时重新生成。
		// 当 MaxAge 设置为 0 时，表示 session 是持久化的，不会自动过期。
		session.Options.MaxAge = -1
		// 重定向
		nextUrl := r.FormValue("next")
		_ = session.Save(r, w)
		if nextUrl == "" {
			nextUrl = r.Referer()
		}
		http.Redirect(w, r, nextUrl, 302)
	})
}
