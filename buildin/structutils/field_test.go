package structutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/structs"
	"github.com/stretchr/testify/assert"
)

// TestValue Get the `Field` underlying value
// func (f *Field) Value() any
func TestValue(t *testing.T) {
	type Parent struct {
		Name string `json:"name,omitempty"`
	}
	p1 := &Parent{"111"}

	s := structs.New(p1)
	n, _ := s.Field("Name")

	assert.Equal(t, "111", n.Value())
}
