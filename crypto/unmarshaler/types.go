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

// ICodec defines codec interface.
type ICodec interface {

	// Name returns codec's name.
	Name() string

	// Unmarshal deserializes the config data bytes into
	// the second input parameter.
	Unmarshal([]byte, interface{}) error
}

// RegisterCodec registers codec by its name.
func RegisterCodec(c ICodec) {
	lock.Lock()
	codecMap[c.Name()] = c
	lock.Unlock()
}

// GetCodec returns the codec by name.
func GetCodec(name string) ICodec {
	lock.RLock()
	c := codecMap[name]
	lock.RUnlock()
	return c
}
