package hscan

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// structMap contains the map of struct fields for target structs
// indexed by the struct type.
// 存放结构体的structSpec对象
type structMap struct {
	m sync.Map
}

// 创建一个全局线程安全字典，key是结构体值运行时表示，value是*structSpec
func newStructMap() *structMap {
	return new(structMap)
}

func (s *structMap) get(t reflect.Type) *structSpec {
	// 从字典中取值，根据结构体值运行时表示获得StrctSpec对象
	if v, ok := s.m.Load(t); ok {
		return v.(*structSpec)
	}
	// 如果key不再字典中，则创建一个新的structSpec对象指针
	spec := newStructSpec(t, "redis")
	// 将新创建的structSpec对象指针放到字段中
	s.m.Store(t, spec)
	return spec
}

//------------------------------------------------------------------------------

// 结构体字段类型解码器对象
// structSpec contains the list of all fields in a target struct.
type structSpec struct {
	m map[string]*structField // key为tag
}

func (s *structSpec) set(tag string, sf *structField) {
	s.m[tag] = sf
}

// 创建结构体字段类型解码器对象，设置每个字段的类型解码器
// 	t: 是结构体值运行时表示
// 	fieldTag: 注解的前缀
func newStructSpec(t reflect.Type, fieldTag string) *structSpec {
	out := &structSpec{
		m: make(map[string]*structField),
	}

	// 返回结构体的字段数量
	num := t.NumField()
	for i := 0; i < num; i++ {
		// 返回字段对象
		f := t.Field(i)
		// 根据前缀获得字段注解的值
		tag := f.Tag.Get(fieldTag)
		if tag == "" || tag == "-" {
			continue
		}
		// 取注解中第一个值
		tag = strings.Split(tag, ",")[0]
		if tag == "" {
			continue
		}
		// Use the built-in decoder.
		// 保存注解
		// i: 字段索引
		// fn: 字段对应的解码器
		out.set(tag, &structField{index: i, fn: decoders[f.Type.Kind()]})
	}

	return out
}

// structField represents a single field in a target struct.
type structField struct {
	index int
	fn    decoderFunc
}

type StructValue struct {
	spec  *structSpec   // 结构体描述
	value reflect.Value // 结构体值的运行时表示
}

// 根据tag名称扫描结构体的指定字段，调用字段类型对应的解码器
func (s StructValue) Scan(key string, value string) error {
	// 根据tag名获得structField对象
	field, ok := s.spec.m[key]
	if !ok {
		// 如果字段没有定义解码器
		return nil
	}
	// 调用字段的解码器
	if err := field.fn(s.value.Field(field.index), value); err != nil {
		// 获得结构体类型对象
		t := s.value.Type()
		return fmt.Errorf("cannot scan redis.result %s into struct field %s.%s of type %s, error-%s",
			value, t.Name(), t.Field(field.index).Name, t.Field(field.index).Type, err.Error())
	}
	return nil
}
