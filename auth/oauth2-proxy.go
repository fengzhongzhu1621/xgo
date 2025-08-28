package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/fengzhongzhu1621/xgo/network/nethttp/staticserver"
)

func HandleOauth2() {
	http.HandleFunc("/-/user", func(w http.ResponseWriter, r *http.Request) {
		// 用于解析 URL 编码的查询字符串（query string），并将其转换为一个 url.Values 类型，后者是一个 map[string][]string 类型的别名，表示键值对的集合。
		fullNameMap, _ := url.ParseQuery(r.Header.Get("X-Auth-Request-Fullname"))
		var fullName string
		for k := range fullNameMap {
			fullName = k
			break
		}
		user := &staticserver.UserInfo{
			Email:    r.Header.Get("X-Auth-Request-Email"),
			Name:     fullName,
			NickName: r.Header.Get("X-Auth-Request-User"),
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		data, _ := json.Marshal(user)
		w.Write(data)
	})
}
