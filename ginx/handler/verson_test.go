package handler

import (
	"net/http"
	"testing"

	"github.com/fengzhongzhu1621/xgo/ginx"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	r := ginx.SetupRouter()
	r.GET("/version", Version)

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
