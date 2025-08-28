package compress

import (
	"bytes"
	"testing"
)

func TestGZip(t *testing.T) {
	buffer := []byte("test")
	var err error
	var b []byte
	if b, err = Gzip(buffer); err != nil {
		t.Error(err)
	}

	if IsGZIP(b) == false {
		t.Error("not gzip format data")
	}

	if b, err = UnGZip(b); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, buffer) {
		t.Error("gzip error")
	}
}
