package tests

import (
	"testing"

	"github.com/gookit/goutil/fmtutil"
	"github.com/stretchr/testify/assert"
)

func TestPrettyJSON(t *testing.T) {
	tests := []any{
		map[string]int{"a": 1},
		struct {
			A int `json:"a"`
		}{1},
	}
	want := `{
    "a": 1
}`
	for _, sample := range tests {
		got, err := fmtutil.PrettyJSON(sample)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	}
}
