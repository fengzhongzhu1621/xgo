package hex

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

type Script struct {
	src  string
	hash string
}

func NewScript(src string) *Script {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script{
		src:  src,
		hash: hex.EncodeToString(h.Sum(nil)),
	}
}

func (s *Script) Hash() string {
	return s.hash
}

// FromHex is a wrapper around hex.Decode to support go1.20
func FromHex(dst []byte, src string) (int, error) {
	return hex.Decode(dst, []byte(src))
}

// AppendEncode appends the hex-encoded src to dst and returns the extended
// buffer.
func AppendEncode(dst, src []byte) []byte {
	d := make([]byte, len(src)*2)
	hex.Encode(d, src)
	dst = append(dst, d...)
	return dst
}
