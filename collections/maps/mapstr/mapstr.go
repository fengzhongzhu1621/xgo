package mapstr

import "fmt"

// MapStr is the map[string]interface{} type
type MapStr map[string]interface{}

// GetValue get value by field from map
func (m MapStr) GetValue(field string) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("map data is nil")
	}
	return m[field], nil
}
