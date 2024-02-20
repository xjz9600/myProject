package leetcode

//https://leetcode.cn/problems/sqrtx/description/

//输入：x = 8
//输出：2
//解释：8 的算术平方根是 2.82842..., 由于返回类型是整数，小数部分将被舍去。

func mySqrt(x int) int {
	if x <= 1 {
		return x
	}
	l := 1
	r := x
	var ans int
	for l <= r {
		mid := (l + r) / 2
		if mid*mid < x {
			ans = mid
			l = mid + 1
		} else if mid*mid == x {
			return mid
		} else {
			r = mid - 1
		}
	}
	return ans
}
