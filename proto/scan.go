package proto

import (
	"encoding"
	"fmt"
	"reflect"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin/bytesconv"
)

// Scan parses bytes `b` to `v` with appropriate type.
// 自动识别v的类型，将b值赋值给v .
func Scan(b []byte, v interface{}) error {
	// 判断v的类型
	switch v := v.(type) {
	case nil:
		return fmt.Errorf("scan(nil)")
	case *string:
		// v是字符串指针，则[]bytes转换为字符串
		*v = bytesconv.BytesToString(b)
		return nil
	case *[]byte:
		// v是字节数组指针
		*v = b
		return nil
	case *int:
		// v是整型指针，将字节数组转换为int类型
		var err error
		*v, err = bytesconv.Atoi(b)
		return err
	case *int8:
		// v是整型指针，将字节数组转换为int8类型
		n, err := bytesconv.ParseInt(b, 10, 8)
		if err != nil {
			return err
		}
		*v = int8(n)
		return nil
	case *int16:
		// v是整型指针，将字节数组转换为int16类型
		n, err := bytesconv.ParseInt(b, 10, 16)
		if err != nil {
			return err
		}
		*v = int16(n)
		return nil
	case *int32:
		// v是整型指针，将字节数组转换为int32类型
		n, err := bytesconv.ParseInt(b, 10, 32)
		if err != nil {
			return err
		}
		*v = int32(n)
		return nil
	case *int64:
		// v是整型指针，将字节数组转换为int64类型
		n, err := bytesconv.ParseInt(b, 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *uint:
		n, err := bytesconv.ParseUint(b, 10, 64)
		if err != nil {
			return err
		}
		*v = uint(n)
		return nil
	case *uint8:
		n, err := bytesconv.ParseUint(b, 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(n)
		return nil
	case *uint16:
		n, err := bytesconv.ParseUint(b, 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(n)
		return nil
	case *uint32:
		n, err := bytesconv.ParseUint(b, 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(n)
		return nil
	case *uint64:
		n, err := bytesconv.ParseUint(b, 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *float32:
		n, err := bytesconv.ParseFloat(b, 32)
		if err != nil {
			return err
		}
		*v = float32(n)
		return err
	case *float64:
		var err error
		*v, err = bytesconv.ParseFloat(b, 64)
		return err
	case *bool:
		// v是bool指针，判断字节数组第一个元素是否为1
		*v = len(b) == 1 && b[0] == '1'
		return nil
	case *time.Time:
		// v是时间对象指针，b的格式必须是time.RFC3339Nano
		var err error
		*v, err = time.Parse(time.RFC3339Nano, bytesconv.BytesToString(b))
		return err
	case encoding.BinaryUnmarshaler:
		// v是解码器对象
		return v.UnmarshalBinary(b)
	default:
		return fmt.Errorf(
			"can't unmarshal %T (consider implementing BinaryUnmarshaler)", v)
	}
}

// 自动识别切片slice的类型，将data中的值转换为切片的类型，并追加到切片slice中 .
func ScanSlice(data []string, slice interface{}) error {
	// 获取数据的运行时表示
	v := reflect.ValueOf(slice)
	// 判断值是否有效。 当值本身非法时，返回 false，例如 reflect Value不包含任何值，值为 nil 等。
	if !v.IsValid() {
		return fmt.Errorf("ScanSlice(nil)")
	}
	// 判断slice是否是一个指针
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("ScanSlice(non-pointer %T)", slice)
	}
	// 获得slice的值，判断是否为切片
	// v是切片指针&[]，v.Elem()是切片对象
	v = v.Elem()
	// v是一个reflect.Value对象，判断v是否是一个切片
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("ScanSlice(non-slice %T)", slice)
	}
	// 返回切片的迭代器
	// v是切片对象
	next := makeSliceNextElemFunc(v)
	for i, s := range data {
		// 获得切片的下一个元素，如果元素是一个指针，则返回指针的值
		elem := next()
		// 自动识别elem的指针的类型，将s值赋值给elem，将切片中每一个元素进行赋值替换
		if err := Scan([]byte(s), elem.Addr().Interface()); err != nil {
			err = fmt.Errorf("ScanSlice index=%d value=%q failed: %w", i, s, err)
			return err
		}
	}

	return nil
}

// 返回一个迭代器，迭代器执行时返回切片的下一个元素，如果元素是一个指针，则返回指针的值
// 	v：是一个切片对象，例如[]*testScanSliceStruct，或[]testScanSliceStruct
func makeSliceNextElemFunc(v reflect.Value) func() reflect.Value {
	// 获得切片的类型
	// v.Type()是[]proto.testScanSliceStruct 或 []proto.testScanSliceStruct
	// v.Type().Elem()是切片中元素的类型对象 proto.testScanSliceStruct 或 *proto.testScanSliceStruct
	elemType := v.Type().Elem()

	// 如果切片中元素的类型是一个指针
	if elemType.Kind() == reflect.Ptr {
		// v =  []
		// v.Type() =  []*proto.testScanSliceStruct
		// elemType =  *proto.testScanSliceStruct
		// elemType.Kind() =  ptr
		elemType = elemType.Elem()
		// elemType是proto.testScanSliceStruct
		return func() reflect.Value {
			if v.Len() < v.Cap() {
				// 切片添加新的一个元素，长度加1
				v.Set(v.Slice(0, v.Len()+1))
				// 获得切片最后一个元素，即新添加的元素
				elem := v.Index(v.Len() - 1)
				// 给新添加的元素设置初始值
				if elem.IsNil() {
					// 创建一个elemType对象，将对象的指针放到切片最后一个元素
					elem.Set(reflect.New(elemType))
				}
				// 返回切片最后一个元素的值
				return elem.Elem()
			}
			// 在切片后面追加一个新的元素
			elem := reflect.New(elemType)
			v.Set(reflect.Append(v, elem))
			// 返回新追加元素的值
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
			// 切片添加新的一个元素，长度加1
			v.Set(v.Slice(0, v.Len()+1))
			// 获得切片最后一个元素，即新添加的元素
			return v.Index(v.Len() - 1)
		}
		// 在切片后面追加一个新的元素
		v.Set(reflect.Append(v, zero))
		// 返回新追加元素的值
		return v.Index(v.Len() - 1)
	}
}
