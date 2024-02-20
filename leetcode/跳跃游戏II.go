package leetcode

//https://leetcode.cn/problems/jump-game-ii/description/

//输入: nums = [2,3,1,1,4]
//输出: 2
//解释: 跳到最后一个位置的最小跳跃数是 2。
//从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。

func jump(nums []int) int {
	var mostRight int
	var step int
	var end int
	for i := 0; i < len(nums)-1; i++ {
		mostRight = max(mostRight, nums[i]+i)
		if i == end {
			end = mostRight
			step++
		}
	}
	return step
}

//func max(a, b int) int {
//	if a >= b {
//		return a
//	}
//	return b
//}
