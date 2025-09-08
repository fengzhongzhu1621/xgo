package plugin

import (
	"errors"

	yaml "gopkg.in/yaml.v3"
)

var _ IDecoder = (*YamlNodeDecoder)(nil)

// YamlNodeDecoder is a decoder for a yaml.Node of the yaml config file.
type YamlNodeDecoder struct {
	Node *yaml.Node
}

// Decode decodes a yaml.Node of the yaml config file.
// 将yaml.Node解码为结构体或字典
func (d *YamlNodeDecoder) Decode(cfg interface{}) error {
	if d.Node == nil {
		return errors.New("yaml node empty")
	}
	return d.Node.Decode(cfg)
}
