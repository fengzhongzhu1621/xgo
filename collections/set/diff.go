package set

import (
	gset "github.com/deckarep/golang-set"
	"github.com/fengzhongzhu1621/xgo/cast"
)

// SetsDifference 将输入的字符串切片转为集合来比较之间的区别
// 接受两个字符串切片 former 和 latter 作为输入，并返回两个接口切片，分别表示在 former 中但不在 latter 中的元素，以及在 latter 中但不在 former 中的元素。
// 此外，如果转换过程中发生错误，它还会返回一个错误。
func SetsDifference(former []string, latter []string) ([]interface{}, []interface{}, error) {
	var (
		err        error
		former_i   []interface{}
		latter_i   []interface{}
		former_set gset.Set
		latter_set gset.Set
	)

	// 将 former 和 latter 字符串切片转换为接口切片
	former_i, err = cast.TransSlice2Interface(former)
	if err != nil {
		return nil, nil, err
	}
	latter_i, err = cast.TransSlice2Interface(latter)
	if err != nil {
		return nil, nil, err
	}

	// 将接口切片转换为集合
	former_set = gset.NewSetFromSlice(former_i)
	latter_set = gset.NewSetFromSlice(latter_i)

	return former_set.Difference(latter_set).ToSlice(), latter_set.Difference(former_set).ToSlice(), nil
}
