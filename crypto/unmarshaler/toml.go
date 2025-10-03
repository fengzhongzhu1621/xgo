package unmarshaler

import "github.com/BurntSushi/toml"

// TomlUnmarshaler is toml unmarshaler.
type TomlUnmarshaler struct{}

// Unmarshal deserializes the data bytes into parameter val in toml protocol.
func (tu *TomlUnmarshaler) Unmarshal(data []byte, val interface{}) error {
	return toml.Unmarshal(data, val)
}

// TomlCodec is toml codec.
type TomlCodec struct{}

// Name returns toml codec's name.
func (*TomlCodec) Name() string {
	return "toml"
}

// Unmarshal deserializes the in bytes into out parameter by toml.
func (c *TomlCodec) Unmarshal(in []byte, out interface{}) error {
	return toml.Unmarshal(in, out)
}
