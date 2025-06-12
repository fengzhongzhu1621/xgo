package buildin

// Go 1.21.0 新增 3 个内置函数，min 和 max 函数，返回 N 个入参中最小/最大的参数，参数类型为 Ordered（有序类型，即支持比较运算符的类型）。
//
// The max built-in function returns the largest value of a fixed number of
// arguments of [cmp.Ordered] types. There must be at least one argument.
// If T is a floating-point type and any of the arguments are NaNs,
// max will return NaN.
// func max[T cmp.Ordered](x T, y ...T) T
