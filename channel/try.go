package channel

// Try 调用函数 func() error，忽略 panic，返回执行是否成功
// calls the function and return false in case of error.
func Try(callback func() error) (ok bool) {
	ok = true

	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	err := callback()
	if err != nil {
		ok = false
	}

	return
}

// Try0 has the same behavior as Try, but callback returns no variable.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try0(callback func()) bool {
	return Try(func() error {
		callback()
		return nil
	})
}

// Try1 is an alias to Try.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try1(callback func() error) bool {
	return Try(callback)
}

// Try2 has the same behavior as Try, but callback returns 2 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try2[T any](callback func() (T, error)) bool {
	return Try(func() error {
		_, err := callback()
		return err
	})
}

// Try3 has the same behavior as Try, but callback returns 3 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try3[T, R any](callback func() (T, R, error)) bool {
	return Try(func() error {
		_, _, err := callback()
		return err
	})
}

// Try4 has the same behavior as Try, but callback returns 4 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try4[T, R, S any](callback func() (T, R, S, error)) bool {
	return Try(func() error {
		_, _, _, err := callback()
		return err
	})
}

// Try5 has the same behavior as Try, but callback returns 5 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try5[T, R, S, Q any](callback func() (T, R, S, Q, error)) bool {
	return Try(func() error {
		_, _, _, _, err := callback()
		return err
	})
}

// Try6 has the same behavior as Try, but callback returns 6 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try6[T, R, S, Q, U any](callback func() (T, R, S, Q, U, error)) bool {
	return Try(func() error {
		_, _, _, _, _, err := callback()
		return err
	})
}

// TryOr has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr[A any](callback func() (A, error), fallbackA A) (A, bool) {
	return TryOr1(callback, fallbackA)
}

// TryOr1 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr1[A any](callback func() (A, error), fallbackA A) (A, bool) {
	ok := false

	Try0(func() {
		a, err := callback()
		if err == nil {
			fallbackA = a
			ok = true
		}
	})

	// 执行成功：返回回调执行结果 和 成功状态
	// 执行失败：返回默认值 和 失败状态
	return fallbackA, ok
}

// TryOr2 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr2[A, B any](callback func() (A, B, error), fallbackA A, fallbackB B) (A, B, bool) {
	ok := false

	Try0(func() {
		a, b, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			ok = true
		}
	})

	return fallbackA, fallbackB, ok
}

// TryOr3 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr3[A, B, C any](
	callback func() (A, B, C, error),
	fallbackA A,
	fallbackB B,
	fallbackC C,
) (A, B, C, bool) {
	ok := false

	Try0(func() {
		a, b, c, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, ok
}

// TryOr4 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr4[A, B, C, D any](
	callback func() (A, B, C, D, error),
	fallbackA A,
	fallbackB B,
	fallbackC C,
	fallbackD D,
) (A, B, C, D, bool) {
	ok := false

	Try0(func() {
		a, b, c, d, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			fallbackD = d
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, fallbackD, ok
}

// TryOr5 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr5[A, B, C, D, E any](
	callback func() (A, B, C, D, E, error),
	fallbackA A,
	fallbackB B,
	fallbackC C,
	fallbackD D,
	fallbackE E,
) (A, B, C, D, E, bool) {
	ok := false

	Try0(func() {
		a, b, c, d, e, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			fallbackD = d
			fallbackE = e
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, fallbackD, fallbackE, ok
}

// TryOr6 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr6[A, B, C, D, E, F any](
	callback func() (A, B, C, D, E, F, error),
	fallbackA A,
	fallbackB B,
	fallbackC C,
	fallbackD D,
	fallbackE E,
	fallbackF F,
) (A, B, C, D, E, F, bool) {
	ok := false

	Try0(func() {
		a, b, c, d, e, f, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			fallbackD = d
			fallbackE = e
			fallbackF = f
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, fallbackD, fallbackE, fallbackF, ok
}

// TryWithErrorValue has the same behavior as Try, but also returns value passed to panic.
// 返回 recover 错误和执行状态
// Play: https://go.dev/play/p/Kc7afQIT2Fs
func TryWithErrorValue(callback func() error) (errorValue any, ok bool) {
	ok = true

	defer func() {
		if r := recover(); r != nil {
			ok = false
			errorValue = r
		}
	}()

	err := callback()
	if err != nil {
		ok = false
		errorValue = err
	}

	return
}

// TryCatch has the same behavior as Try, but calls the catch function in case of error.
// Play: https://go.dev/play/p/PnOON-EqBiU
func TryCatch(callback func() error, catch func()) {
	if !Try(callback) {
		catch()
	}
}

// TryCatchWithErrorValue has the same behavior as TryWithErrorValue, but calls the catch function in case of error.
// Play: https://go.dev/play/p/8Pc9gwX_GZO
func TryCatchWithErrorValue(callback func() error, catch func(any)) {
	if err, ok := TryWithErrorValue(callback); !ok {
		catch(err)
	}
}
