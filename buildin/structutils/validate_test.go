package structutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/structs"
	"github.com/stretchr/testify/assert"
)

// TesIsStruct Check if the struct is valid
// func (s *Struct) IsStruct() bool
func TestIsStruct(t *testing.T) {
	type People struct {
		Name string `json:"name"`
	}
	p1 := &People{Name: "11"}
	s := structs.New(p1)

	assert.Equal(t, true, s.IsStruct())
}
