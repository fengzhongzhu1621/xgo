package reflectutils

import "reflect"

// TypeConverter 类型转换器
type TypeConverter struct {
	SrcType interface{}
	DstType interface{}
	Fn      func(src interface{}) (dst interface{}, err error)
}

type ConverterPair struct {
	SrcType reflect.Type
	DstType reflect.Type
}

// Option sets copy options 转换器配置项
type Option struct {
	// setting this value to true will ignore copying zero values of all the fields, including bools, as well as a
	// struct having all it's fields set to their zero values respectively (see IsZero() in reflect/value.go)
	IgnoreEmpty bool
	DeepCopy    bool
	Converters  []TypeConverter // 包含多个类型转化器
}

// TypeConverters 将类型转换器转换为具有反射类型的字典
func (opt Option) TypeConverters() map[ConverterPair]TypeConverter {
	converters := map[ConverterPair]TypeConverter{}

	// save converters into map for faster lookup
	for i := range opt.Converters {
		pair := ConverterPair{
			SrcType: reflect.TypeOf(opt.Converters[i].SrcType),
			DstType: reflect.TypeOf(opt.Converters[i].DstType),
		}

		converters[pair] = opt.Converters[i]
	}

	return converters
}
