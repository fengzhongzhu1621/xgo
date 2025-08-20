package math

// PowerToRoundUp 计算能够覆盖给定整数 n 的最小 2 的幂的指数值。换句话说，它找到最小的 power 使得 2^power >= n。
func PowerToRoundUp(n int) int {
	powerOfTwo, power := 1, 0         // 初始化：2^0 = 1，指数为0
	for ; n-powerOfTwo > 0; power++ { // 循环条件：当n > powerOfTwo时继续
		powerOfTwo <<= 1 // 每次循环将powerOfTwo左移1位（即乘以2）
	}
	return power // 返回最终的指数值
}
