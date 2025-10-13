package entity

// Entity 配置实体，包含配置内容的原始数据和解析后的数据
type Entity struct {
	// 配置原始内容
	raw []byte // current binary data
	// 解析后的配置对象
	data interface{} // unmarshal type to use point type, save latest no error data
}

func NewEntity() *Entity {
	return &Entity{
		raw:  []byte{},
		data: make(map[string]interface{}),
	}
}

func (e *Entity) GetRaw() []byte {
	return e.raw
}

func (e *Entity) GetData() interface{} {
	return e.data
}

func (e *Entity) SetRaw(raw []byte) {
	e.raw = raw
}

func (e *Entity) SetData(data interface{}) {
	e.data = data
}
