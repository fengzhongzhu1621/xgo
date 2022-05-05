package hcl

// HCL 是 hashicorp 推出的一个配置语言，在 hashicorp 的产品，
// 如 Consul、Terraform 中用作标准的配置语言，其语法结构有点类似于 nginx 的配置文件 .

import (
	"bytes"
	"encoding/json"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
)

// Codec implements the encoding.Encoder and encoding.Decoder interfaces for HCL encoding.
// TODO: add printer config to the codec?
type Codec struct{}

func (Codec) Encode(v map[string]interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// TODO: use printer.Format? Is the trailing newline an issue?

	ast, err := hcl.Parse(string(b))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = printer.Fprint(&buf, ast.Node)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (Codec) Decode(b []byte, v map[string]interface{}) error {
	return hcl.Unmarshal(b, &v)
}
