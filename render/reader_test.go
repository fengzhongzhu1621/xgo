package render

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRenderReader(t *testing.T) {
	w := httptest.NewRecorder()

	body := "#!PNG some raw data"
	headers := make(map[string]string)
	headers["Content-Disposition"] = `attachment; filename="filename.png"`
	headers["x-request-id"] = "requestId"

	err := (Reader{
		ContentLength: int64(len(body)),
		ContentType:   "image/png",
		Reader:        strings.NewReader(body),
		Headers:       headers,
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, body, w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.Itoa(len(body)), w.Header().Get("Content-Length"))
	assert.Equal(t, headers["Content-Disposition"], w.Header().Get("Content-Disposition"))
	assert.Equal(t, headers["x-request-id"], w.Header().Get("x-request-id"))
}

func TestRenderReaderNoContentLength(t *testing.T) {
	w := httptest.NewRecorder()

	body := "#!PNG some raw data"
	headers := make(map[string]string)
	headers["Content-Disposition"] = `attachment; filename="filename.png"`
	headers["x-request-id"] = "requestId"

	err := (Reader{
		ContentLength: -1,
		ContentType:   "image/png",
		Reader:        strings.NewReader(body),
		Headers:       headers,
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, body, w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	assert.NotContains(t, "Content-Length", w.Header())
	assert.Equal(t, headers["Content-Disposition"], w.Header().Get("Content-Disposition"))
	assert.Equal(t, headers["x-request-id"], w.Header().Get("x-request-id"))
}

