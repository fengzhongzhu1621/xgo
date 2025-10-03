package unmarshaler

func init() {
	RegisterUnmarshaler("yaml", &YamlUnmarshaler{})
	RegisterUnmarshaler("json", &JSONUnmarshaler{})
	RegisterUnmarshaler("toml", &TomlUnmarshaler{})
}

func init() {
	RegisterCodec(&YamlCodec{})
	RegisterCodec(&JSONCodec{})
	RegisterCodec(&TomlCodec{})
}
