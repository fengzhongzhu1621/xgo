package unmarshaler

import (
	yaml "gopkg.in/yaml.v3"
)

// YamlUnmarshaler is yaml unmarshaler.
type YamlUnmarshaler struct{}

// Unmarshal deserializes the data bytes into parameter val in yaml protocol.
func (yu *YamlUnmarshaler) Unmarshal(data []byte, val interface{}) error {
	return yaml.Unmarshal(data, val)
}

// YamlCodec is yaml codec.
type YamlCodec struct{}

// Name returns yaml codec's name.
func (*YamlCodec) Name() string {
	return "yaml"
}

// Unmarshal deserializes the in bytes into out parameter by yaml.
func (c *YamlCodec) Unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}
