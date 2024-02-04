package leetcode

// 输入：height = [0,1,0,2,1,0,1,3,2,1,2,1]
// 输出：6
// 解释：上面是由数组 [0,1,0,2,1,0,1,3,2,1,2,1] 表示的高度图，在这种情况下，可以接 6 个单位的雨水（蓝色部分表示雨水）

// https://leetcode.cn/problems/trapping-rain-water/solutions/692342/jie-yu-shui-by-leetcode-solution-tuvc/
func trap(height []int) int {
	var stack []int
	var result int
	for i, h := range height {
		for len(stack) > 0 && height[stack[len(stack)-1]] < h {
			low := height[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				result += (minTrap(height[stack[len(stack)-1]], h) - low) * (i - stack[len(stack)-1] - 1)
			}
		}
		stack = append(stack, i)
	}
	return result
}

func minTrap(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
