package buntdb

import "github.com/tidwall/grect"

// rect is used by Intersects and Nearby
// 在一维的情况下，是一条线段，用首尾两个点表示
// 在二维的情况下，我们称之为最小限定矩形。 MBR(minimum bounding retangle)
// 通常，我们只需要两个点就可限定一个矩形，也就是矩形某个对角线的两个点就可以决定一个唯一的矩形。
// 通常我们使用（左下，右上两个点表示）或者使用右上左下，都是可以的。
type rect struct {
	min, max []float64
}

func (r *rect) Rect(ctx interface{}) (min, max []float64) {
	return r.min, r.max
}

// Rect is helper function that returns a string representation
// of a rect. IndexRect() is the reverse function and can be used
// to generate a rect from a string.
// 获得线段的字符串表示
func Rect(min, max []float64) string {
	r := grect.Rect{Min: min, Max: max}
	return r.String()
}

// Point is a helper function that converts a series of float64s
// to a rectangle for a spatial index.
// 将多边形转换为矩形表示
func Point(coords ...float64) string {
	return Rect(coords, coords)
}

// IndexRect is a helper function that converts string to a rect.
// Rect() is the reverse function and can be used to generate a string
// from a rect.
// 将矩形的字符串表示转换为对象
func IndexRect(a string) (min, max []float64) {
	r := grect.Get(a)
	return r.Min, r.Max
}
