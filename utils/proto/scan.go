package proto

import (
	"encoding"
	"fmt"
	"reflect"
	"time"

	"xgo/utils/bytes_utils"
	"xgo/utils/bytesconv"
)

// Scan parses bytes `b` to `v` with appropriate type.
// 自动识别v的类型，将b值赋值给v
func Scan(b []byte, v interface{}) error {
	switch v := v.(type) {
	case nil:
		return fmt.Errorf("redis: Scan(nil)")
	case *string:
		*v = bytesconv.BytesToString(b)
		return nil
	case *[]byte:
		*v = b
		return nil
	case *int:
		var err error
		*v, err = bytes_utils.Atoi(b)
		return err
	case *int8:
		n, err := bytes_utils.ParseInt(b, 10, 8)
		if err != nil {
			return err
		}
		*v = int8(n)
		return nil
	case *int16:
		n, err := bytes_utils.ParseInt(b, 10, 16)
		if err != nil {
			return err
		}
		*v = int16(n)
		return nil
	case *int32:
		n, err := bytes_utils.ParseInt(b, 10, 32)
		if err != nil {
			return err
		}
		*v = int32(n)
		return nil
	case *int64:
		n, err := bytes_utils.ParseInt(b, 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *uint:
		n, err := bytes_utils.ParseUint(b, 10, 64)
		if err != nil {
			return err
		}
		*v = uint(n)
		return nil
	case *uint8:
		n, err := bytes_utils.ParseUint(b, 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(n)
		return nil
	case *uint16:
		n, err := bytes_utils.ParseUint(b, 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(n)
		return nil
	case *uint32:
		n, err := bytes_utils.ParseUint(b, 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(n)
		return nil
	case *uint64:
		n, err := bytes_utils.ParseUint(b, 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *float32:
		n, err := bytes_utils.ParseFloat(b, 32)
		if err != nil {
			return err
		}
		*v = float32(n)
		return err
	case *float64:
		var err error
		*v, err = bytes_utils.ParseFloat(b, 64)
		return err
	case *bool:
		*v = len(b) == 1 && b[0] == '1'
		return nil
	case *time.Time:
		var err error
		*v, err = time.Parse(time.RFC3339Nano, bytesconv.BytesToString(b))
		return err
	case encoding.BinaryUnmarshaler:
		return v.UnmarshalBinary(b)
	default:
		return fmt.Errorf(
			"redis: can't unmarshal %T (consider implementing BinaryUnmarshaler)", v)
	}
}

// 自动识别切片slice的类型，将data中的值转换为切片的类型，并追加到切片slice中
func ScanSlice(data []string, slice interface{}) error {
	// 获取数据的运行时表示
	v := reflect.ValueOf(slice)
	// 判断值是否有效。 当值本身非法时，返回 false，例如 reflect Value不包含任何值，值为 nil 等。
	if !v.IsValid() {
		return fmt.Errorf("redis: ScanSlice(nil)")
	}
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("redis: ScanSlice(non-pointer %T)", slice)
	}
	// 获得slice的值，判断是否为切片
	// v是切片指针&[]，v.Elem()是切片对象
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("redis: ScanSlice(non-slice %T)", slice)
	}
	// 返回切片的迭代器
	// v是切片对象
	next := makeSliceNextElemFunc(v)
	for i, s := range data {
		// 获得切片的下一个元素
		elem := next()
		// 自动识别elem的指针的类型，将s值赋值给elem
		if err := Scan([]byte(s), elem.Addr().Interface()); err != nil {
			err = fmt.Errorf("redis: ScanSlice index=%d value=%q failed: %w", i, s, err)
			return err
		}
	}

	return nil
}

// 返回一个迭代器
// v：是一个切片对象
func makeSliceNextElemFunc(v reflect.Value) func() reflect.Value {
	// 获得切片的类型
	// v.Type()是[]proto.testScanSliceStruct
	// v.Type().Elem()是切片的类型对象 proto.testScanSliceStruct
	elemType := v.Type().Elem()

	// 切片类型是一个指针
	if elemType.Kind() == reflect.Ptr {
		// v =  []
		// v.Type() =  []*proto.testScanSliceStruct
		// elemType =  *proto.testScanSliceStruct
		// elemType.Kind() =  ptr
		elemType = elemType.Elem()
		return func() reflect.Value {
			if v.Len() < v.Cap() {
				v.Set(v.Slice(0, v.Len()+1))
				elem := v.Index(v.Len() - 1)
				if elem.IsNil() {
					elem.Set(reflect.New(elemType))
				}
				return elem.Elem()
			}

			elem := reflect.New(elemType)
			v.Set(reflect.Append(v, elem))
			return elem.Elem()
		}
	}

	// 切片类型是一个结构体
	// 获得类型的0值
	// v =  []
	// v.Type() =  []proto.testScanSliceStruct
	// elemType =  proto.testScanSliceStruct
	// elemType.Kind() =  struct
	// zero =  {0 }
	zero := reflect.Zero(elemType)
	return func() reflect.Value {
		if v.Len() < v.Cap() {
			// 提取字符串
			v.Set(v.Slice(0, v.Len()+1))
			// 返回最后一条数据
			return v.Index(v.Len() - 1)
		}

		v.Set(reflect.Append(v, zero))
		return v.Index(v.Len() - 1)
	}
}
