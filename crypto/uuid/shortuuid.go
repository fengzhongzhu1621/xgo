package uuid

import (
	"github.com/lithammer/shortuuid/v3"
)

// NewShortUUID returns a new short UUID.
func NewShortUUID() string {
	return shortuuid.New()
}
