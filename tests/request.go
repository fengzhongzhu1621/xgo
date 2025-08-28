package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// CreateTestContextWithDefaultRequest ...
func CreateTestContextWithDefaultRequest(w *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest("POST", "/", new(bytes.Buffer))
	return ctx
}
