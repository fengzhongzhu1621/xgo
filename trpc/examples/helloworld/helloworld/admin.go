package main

import (
	"net/http"
)

func customAdminCmd(w http.ResponseWriter, r *http.Request) {
	// 业务逻辑

	// 返回成功错误码
	w.Write([]byte(`{"errorcode":0, "message":"ok"}`))
}
