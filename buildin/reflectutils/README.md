# reflecdt

## reflect.TypeOf
返回一个reflect.Type类型的值，表示给定值的类型。


## reflect.ValueOf 返回给定值的值


## reflect.Value 表示运行时数据的结构类型


## reflect.Type 表示运行时类型信息的结构类型 int, string, float, struct


## reflect.Kind 类型的枚举 reflect.Int, reflect.Struct


## reflect.Value.Method 返回与给定名称的方法对应的函数值


# reflect.DeepEqual
比较逻辑
* 基本类型：直接比较值是否相等。
* 数组：比较每个元素是否相等。
* 切片：比较切片长度和每个元素是否相等。
* 映射：比较键值对的数量和每个键对应的值是否相等。
* 结构体：比较每个字段是否相等。
* 指针：比较指针指向的值是否相等。
* 接口：比较接口的动态类型和值是否相等。

局限性
* 未导出字段：对于结构体的未导出字段，reflect.DeepEqual无法访问，因此无法比较。
* 函数比较：如果结构体或容器中包含函数，reflect.DeepEqual不会比较函数逻辑是否相同，只会比较函数指针是否相等。
