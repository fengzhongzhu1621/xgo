package unmarshaler

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
