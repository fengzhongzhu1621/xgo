package handler

import (
	"net/http"
	"testing"

	"github.com/fengzhongzhu1621/xgo/ginx"
	"github.com/steinfletcher/apitest"
)

func TestPing(t *testing.T) {
	t.Parallel()

	// 注册路由
	r := ginx.SetupRouter()
	r.GET("/ping", Ping)

	// 模拟调用
	apitest.New().
		Handler(r).
		Get("/ping").
		Expect(t).
		Body(`pong`).
		Status(http.StatusOK).
		End()
}
