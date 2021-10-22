package hscan

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// 定义结构体注解，提供Scan()给结构体的字段赋值，每个字段赋值前进行扫描校验，
// 将字符串转换为指定类型，如果有一个字段的赋值不符合策略，则停止赋值，返回err
// decoderFunc represents decoding functions for default built-in types.
type decoderFunc func(reflect.Value, string) error

var (
	// List of built-in decoders indexed by their numeric constant values (eg: reflect.Bool = 1).
	// 定义切片字典
	decoders = []decoderFunc{
		reflect.Bool:          decodeBool,
		reflect.Int:           decodeInt,
		reflect.Int8:          decodeInt8,
		reflect.Int16:         decodeInt16,
		reflect.Int32:         decodeInt32,
		reflect.Int64:         decodeInt64,
		reflect.Uint:          decodeUint,
		reflect.Uint8:         decodeUint8,
		reflect.Uint16:        decodeUint16,
		reflect.Uint32:        decodeUint32,
		reflect.Uint64:        decodeUint64,
		reflect.Float32:       decodeFloat32,
		reflect.Float64:       decodeFloat64,
		reflect.Complex64:     decodeUnsupported,
		reflect.Complex128:    decodeUnsupported,
		reflect.Array:         decodeUnsupported,
		reflect.Chan:          decodeUnsupported,
		reflect.Func:          decodeUnsupported,
		reflect.Interface:     decodeUnsupported,
		reflect.Map:           decodeUnsupported,
		reflect.Ptr:           decodeUnsupported,
		reflect.Slice:         decodeSlice,
		reflect.String:        decodeString,
		reflect.Struct:        decodeUnsupported,
		reflect.UnsafePointer: decodeUnsupported,
	}

	// Global map of struct field specs that is populated once for every new
	// struct type that is scanned. This caches the field types and the corresponding
	// decoder functions to avoid iterating through struct fields on subsequent scans.
	// 定义全局字典（线程安全），包含结构体中的structSpec对象
	globalStructMap = newStructMap()
)

// 将结构体指针转换为StructValue对象，dst是结构体指针
func Struct(dst interface{}) (StructValue, error) {
	// 获取结构体指针的运行时表示
	v := reflect.ValueOf(dst)
	// 判断dst是否是一个结构体指针
	// The destination to scan into should be a struct pointer.
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return StructValue{}, fmt.Errorf("redis.Scan(non-pointer %T)", dst)
	}
	// 因为v是一个指针，获取对应值的运行时表示，及结构体对象运行时表示
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return StructValue{}, fmt.Errorf("redis.Scan(non-struct %T)", dst)
	}
	// 解析结构体字段的注解，获得结构指针的值和所有字段的类型解码器对象
	return StructValue{
		spec:  globalStructMap.get(v.Type()), // 获得结构体描述，包含每个字段的类型解码器
		value: v,                             // 结构体的运行时表示
	}, nil
}

// Scan scans the results from a key-value Redis map result set to a destination struct.
// The Redis keys are matched to the struct's field with the `redis` tag.
// 扫描结构体的所有字段，调用字段类型对应的解码器
// 在给结构体的字段赋值前，进行扫描校验，将字符串转换为指定类型，如果有一个字段的赋值不符合策略，则停止赋值，返回err
// Args：
// 	ds: 结构体指针
// 	keys: 结构体中字段的tag名数组
//  vals: 需要给结构体赋值的数组
func Scan(dst interface{}, keys []interface{}, vals []interface{}) error {
	if len(keys) != len(vals) {
		return errors.New("args should have the same number of keys and vals")
	}
	// 将结构体指针转换为StructValue对象
	strct, err := Struct(dst)
	if err != nil {
		return err
	}

	// Iterate through the (key, value) sequence.
	for i := 0; i < len(vals); i++ {
		// 获得tag名
		key, ok := keys[i].(string)
		if !ok {
			continue
		}
		// 获得需要设置的字段的值
		val, ok := vals[i].(string)
		if !ok {
			continue
		}
		// 根据tag名称扫描结构体的指定字段，调用字段类型对应的解码器，给结构体字段赋值
		if err := strct.Scan(key, val); err != nil {
			return err
		}
	}

	return nil
}

func decodeBool(f reflect.Value, s string) error {
	// 将字符串转换为bool类型
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	f.SetBool(b)
	return nil
}

func decodeInt8(f reflect.Value, s string) error {
	return decodeNumber(f, s, 8)
}

func decodeInt16(f reflect.Value, s string) error {
	return decodeNumber(f, s, 16)
}

func decodeInt32(f reflect.Value, s string) error {
	return decodeNumber(f, s, 32)
}

func decodeInt64(f reflect.Value, s string) error {
	return decodeNumber(f, s, 64)
}

func decodeInt(f reflect.Value, s string) error {
	return decodeNumber(f, s, 0)
}

func decodeNumber(f reflect.Value, s string, bitSize int) error {
	v, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		return err
	}
	f.SetInt(v)
	return nil
}

func decodeUint8(f reflect.Value, s string) error {
	return decodeUnsignedNumber(f, s, 8)
}

func decodeUint16(f reflect.Value, s string) error {
	return decodeUnsignedNumber(f, s, 16)
}

func decodeUint32(f reflect.Value, s string) error {
	return decodeUnsignedNumber(f, s, 32)
}

func decodeUint64(f reflect.Value, s string) error {
	return decodeUnsignedNumber(f, s, 64)
}

func decodeUint(f reflect.Value, s string) error {
	return decodeUnsignedNumber(f, s, 0)
}

func decodeUnsignedNumber(f reflect.Value, s string, bitSize int) error {
	v, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		return err
	}
	f.SetUint(v)
	return nil
}

func decodeFloat32(f reflect.Value, s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	f.SetFloat(v)
	return nil
}

// although the default is float64, but we better define it.
func decodeFloat64(f reflect.Value, s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	f.SetFloat(v)
	return nil
}

func decodeString(f reflect.Value, s string) error {
	f.SetString(s)
	return nil
}

func decodeSlice(f reflect.Value, s string) error {
	// []byte slice ([]uint8).
	if f.Type().Elem().Kind() == reflect.Uint8 {
		f.SetBytes([]byte(s))
	}
	return nil
}

func decodeUnsupported(v reflect.Value, s string) error {
	return fmt.Errorf("redis.Scan(unsupported %s)", v.Type())
}
