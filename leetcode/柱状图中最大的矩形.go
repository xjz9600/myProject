package leetcode

// https://leetcode.cn/problems/0ynMMM/submissions/497721334/
//输入：heights = [2,1,5,6,2,3]
//输出：10
//解释：最大的矩形为图中红色区域，面积为 10
//提示：记录左右两边的位置，使用单调栈

func largestRectangleArea(heights []int) int {
	left := make([]int, len(heights))
	right := make([]int, len(heights))
	var stack []int
	for i := 0; i < len(heights); i++ {
		right[i] = len(heights)
	}
	for i := 0; i < len(heights); i++ {
		for len(stack) > 0 && heights[stack[len(stack)-1]] >= heights[i] {
			right[stack[len(stack)-1]] = i
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			left[i] = -1
		} else {
			left[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	var ans int
	for i := 0; i < len(left); i++ {
		ans = maxAns(ans, (right[i]-left[i]-1)*heights[i])
	}
	return ans
}

func maxAns(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
