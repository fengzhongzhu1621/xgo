package unmarshaler

import "encoding/json"

// JSONUnmarshaler is json unmarshaler.
type JSONUnmarshaler struct{}

// Unmarshal deserializes the data bytes into parameter val in json protocol.
func (ju *JSONUnmarshaler) Unmarshal(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

// JSONCodec is json codec.
type JSONCodec struct{}

// Name returns json codec's name.
func (*JSONCodec) Name() string {
	return "json"
}

// Unmarshal deserializes the in bytes into out parameter by json.
func (c *JSONCodec) Unmarshal(in []byte, out interface{}) error {
	return json.Unmarshal(in, out)
}
