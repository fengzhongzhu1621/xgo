package file

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/stretchr/testify/assert"
)

// Get file mime type, 'file' param's type should be string or *os.File.
// func MiMeType(file any) string
func TestMiMeType(t *testing.T) {
	fname := "./test.txt"
	file, _ := os.Create(fname)
	file.WriteString("hello world")

	f, _ := os.Open(fname)
	defer f.Close()

	mimeType := fileutil.MiMeType(f)
	fmt.Println(mimeType)

	os.Remove(fname)

	// Output:
	// application/octet-stream
}

func TestMimeType(t *testing.T) {
	assert.Equal(t, "", fsutil.DetectMime(""))
	assert.Equal(t, "", fsutil.MimeType("not-exist"))
	assert.Equal(t, "image/jpeg", fsutil.MimeType("testdata/test.jpg"))

	buf := new(bytes.Buffer)
	buf.Write([]byte("\xFF\xD8\xFF"))
	assert.Equal(t, "image/jpeg", fsutil.ReaderMimeType(buf))
	buf.Reset()

	buf.Write([]byte("text"))
	assert.Equal(t, "text/plain; charset=utf-8", fsutil.ReaderMimeType(buf))
	buf.Reset()

	buf.Write([]byte(""))
	assert.Equal(t, "", fsutil.ReaderMimeType(buf))
	buf.Reset()

	assert.True(t, fsutil.IsImageFile("testdata/test.jpg"))
	assert.False(t, fsutil.IsImageFile("testdata/not-exists"))
}
