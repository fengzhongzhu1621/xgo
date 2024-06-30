package crypto

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
