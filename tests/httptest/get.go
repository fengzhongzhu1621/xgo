package httptest

import (
	"io"
	"net/http"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func HttpGet(url string, expected1 string, expected2 string, t *testing.T) {
	resp, err := http.Get(url)
	assert.Nil(t, err)
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Nil(t, err)
	if string(body) != expected1 && string(body) != expected2 {
		t.Errorf("Expected %#v or %#v, got %#v.", expected1, expected2, string(body))
	}
}
