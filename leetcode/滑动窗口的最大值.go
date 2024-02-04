package leetcode

// https://leetcode.cn/problems/sliding-window-maximum/submissions/497939288/
// 输入：nums = [1,3,-1,-3,5,3,6,7], k = 3
// 输出：[3,3,5,5,6,7]
// 解释：
// 滑动窗口的位置                最大值
// ---------------               -----
// [1  3  -1] -3  5  3  6  7       3
// 1 [3  -1  -3] 5  3  6  7       3
// 1  3 [-1  -3  5] 3  6  7       5
// 1  3  -1 [-3  5  3] 6  7       5
// 1  3  -1  -3 [5  3  6] 7       6
// 1  3  -1  -3  5 [3  6  7]      7

// 提示一定要记录位置否则[7,2,4,1] k=3 的时候就出现问题
func maxSlidingWindow(nums []int, k int) []int {
	var stack []int
	var push = func(k int) {
		for len(stack) > 0 && nums[stack[len(stack)-1]] <= nums[k] {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, k)
	}
	for i := 0; i < k; i++ {
		push(i)
	}
	var ans []int
	ans = append(ans, nums[stack[0]])
	for i := k; i < len(nums); i++ {
		push(i)
		for stack[0] <= i-k {
			stack = stack[1:]
		}
		ans = append(ans, nums[stack[0]])
	}
	return ans
}
