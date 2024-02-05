package leetcode

// https://leetcode.cn/problems/powx-n/submissions/500472800/
// 实现 pow(x, n) ，即计算 x 的整数 n 次幂函数（即，xn ）。
func myPow(x float64, n int) float64 {
	var result float64
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	if n == -1 {
		return 1 / x
	}
	if n > 0 {
		result = myPow(x, n/2)
		result = result * result
		if n%2 == 1 {
			result = result * x
		}
	} else {
		result = myPow(x, n/2)
		result = result * result
		if n%2 == -1 {
			result = result * 1 / x
		}
	}
	return result
}
