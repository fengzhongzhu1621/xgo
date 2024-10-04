package handler_test

import (
	"net/http"
	"testing"

	"github.com/fengzhongzhu1621/xgo/ginx"
	"github.com/fengzhongzhu1621/xgo/ginx/handler"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Parallel()

	// 注册路由
	r := ginx.SetupRouter()
	r.GET("/ping", handler.Ping)

	// 模拟调用
	apitest.New().
		Handler(r).
		Get("/ping").
		Expect(t).
		Body(`{"message":"pong"}`).
		Status(http.StatusOK).
		End()
}

func TestVersion(t *testing.T) {
	t.Parallel()

	r := ginx.SetupRouter()
	r.GET("/version", handler.Version)

	apitest.New().
		Handler(r).
		Get("/version").
		Expect(t).
		// 校验响应结果
		Assert(ginx.NewJSONAssertFunc(t, func(m map[string]interface{}) error {
			assert.Contains(t, m, "version")
			assert.Contains(t, m, "commit")
			assert.Contains(t, m, "buildTime")
			assert.Contains(t, m, "goVersion")
			assert.Contains(t, m, "env")
			return nil
		})).
		Status(http.StatusOK).
		End()
}
