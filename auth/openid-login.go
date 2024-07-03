package auth

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/codeskyblue/openid-go"
	"github.com/gorilla/sessions"
)

var (
	// 创建一个新的SimpleNonceStore实例。SimpleNonceStore是一个用于存储和检索随机数（nonce）的简单容器
	// 在密码学中，nonce 是一个只使用一次的数值，通常与密码散列函数一起使用，以确保先前的散列不会在未来被重复使用。这对于防止重放攻击特别重要。
	nonceStore = openid.NewSimpleNonceStore()
	// 用于创建一个新的SimpleDiscoveryCache实例。SimpleDiscoveryCache是一个用于缓存OpenID提供者元数据的简单容器。
	// 这些元数据包括提供者的端点URL、公钥等，用于在OpenID认证过程中进行验证和授权。
	discoveryCache = openid.NewSimpleDiscoveryCache()
	// 用于创建一个新的基于 cookie 的 session 存储引擎。
	// 接受一个密钥（key）切片作为参数，这些密钥用于对 cookie 数据进行签名和加密，以确保数据的安全性和完整性。
	store              = sessions.NewCookieStore([]byte("something-very-secret"))
	defaultSessionName = "ghs-session"
)

type UserInfo struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	NickName string `json:"nickName"`
}

type M map[string]interface{}

func init() {
	// 注册用户自定义的类型，以便在编码和解码时使用
	gob.Register(&UserInfo{})
	gob.Register(&M{})
}

// 注册 openid 登录相关的路由
func handleOpenID(loginUrl string, secure bool) {
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
		id, err := openid.Verify("http://"+r.Host+r.URL.String(), discoveryCache, nonceStore)
		if err != nil {
			io.WriteString(w, "Authentication check failed.")
			return
		}
		// 将用户信息存放到 session 中
		session, err := store.Get(r, defaultSessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := &UserInfo{
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

	http.HandleFunc("/-/user", func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, defaultSessionName)
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
		session, err := store.Get(r, defaultSessionName)
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
