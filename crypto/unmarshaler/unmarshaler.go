package unmarshaler

import "sync"

var (
	codecMap = make(map[string]ICodec)
	lock     = sync.RWMutex{}
)
var (
	unmarshalers = make(map[string]IUnmarshaler)
)

// IUnmarshaler defines a unmarshal interface, this will
// be used to parse config data.
type IUnmarshaler interface {
	// Unmarshal deserializes the data bytes into value parameter.
	Unmarshal(data []byte, value interface{}) error
}

// RegisterUnmarshaler registers an unmarshaler by name.
func RegisterUnmarshaler(name string, us IUnmarshaler) {
	unmarshalers[name] = us
}

// GetUnmarshaler returns an unmarshaler by name.
func GetUnmarshaler(name string) IUnmarshaler {
	return unmarshalers[name]
}
