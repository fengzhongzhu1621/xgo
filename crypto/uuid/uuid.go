package uuid

import (
	"crypto/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

func New() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// NewULID returns a new ULID.
func NewULID() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}

// GenerateID 获得随机字符串.
func GenerateID() string {
	base := 10
	return strconv.FormatInt(time.Now().UnixNano(), base)
}
