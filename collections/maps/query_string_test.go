package maps

import (
	"fmt"
	"testing"

	goutil_maputil "github.com/gookit/goutil/maputil"
	"github.com/stretchr/testify/assert"
)

func TestHTTPQueryString(t *testing.T) {
	src := map[string]any{"a": "v0", "b": 23}
	str := goutil_maputil.HTTPQueryString(src)

	fmt.Println(str)
	assert.Contains(t, str, "b=23")
	assert.Contains(t, str, "a=v0")
}
