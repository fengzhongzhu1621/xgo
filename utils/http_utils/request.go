package http_utils

import (
	"bytes"
	"net/http"
)

func RequestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
