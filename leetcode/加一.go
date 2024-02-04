package leetcode

// https://leetcode.cn/problems/plus-one/submissions/499616303/

//输入：digits = [1,2,3]
//输出：[1,2,4]
//解释：输入数组表示数字 123。

//输入：digits = [9,9,9]
//输出：[1,0,0,0]
//解释：输入数组表示数字 999。

func plusOne(digits []int) []int {
	var j = 1
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i]+j > 9 {
			digits[i] = digits[i] + j - 10
			j = 1
		} else {
			digits[i] = digits[i] + j
			return digits
		}
	}
	var result = make([]int, len(digits)+1)
	result[0] = 1
	return result
}
