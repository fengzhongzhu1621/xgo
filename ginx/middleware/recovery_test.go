package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestPanicLog(t *testing.T) {
	t.Parallel()

	err1 := errors.New("test")
	panicLog(err1)
}

func TestIsBrokenPipeError(t *testing.T) {
	t.Parallel()

	err1 := errors.New("test")
	assert.False(t, isBrokenPipeError(err1))

	err2 := &net.OpError{
		Op:     "",
		Net:    "",
		Source: nil,
		Addr:   nil,
		Err: &os.SyscallError{
			Syscall: "",
			Err:     errors.New("broken pipe"),
		},
	}
	assert.True(t, isBrokenPipeError(err2))
}

func TestRecovery(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.Use(Recovery(true))

	r.GET("/ping", func(c *gin.Context) {
		b := 0
		a := 1 / b
		fmt.Println(a, b)
		c.String(200, "pong")
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}
