package leetcode

//https://leetcode.cn/problems/valid-perfect-square/description/

//输入：num = 16
//输出：true
//解释：返回 true ，因为 4 * 4 = 16 且 4 是一个整数。

func isPerfectSquare(num int) bool {
	if num <= 1 {
		return true
	}
	l := 1
	r := num
	for l <= r {
		mid := (l + r) / 2
		if mid*mid < num {
			l = mid + 1
		} else if mid*mid == num {
			return true
		} else {
			r = mid - 1
		}
	}
	return false
}
