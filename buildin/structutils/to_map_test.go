package structutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/structs"
)

// TestToMap convert a valid struct to a map
// func (s *Struct) ToMap() (map[string]any, error)
// func ToMap(v any) (map[string]any, error)
func TestToMap(t *testing.T) {
	type People struct {
		Name string `json:"name"`
	}
	p1 := &People{Name: "11"}
	// use constructor function
	s1 := structs.New(p1)
	m1, _ := s1.ToMap()
	fmt.Println(m1)

	// use static function
	m2, _ := structs.ToMap(p1)
	fmt.Println(m2)
}
