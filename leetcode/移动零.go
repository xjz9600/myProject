package leetcode

// https://leetcode.cn/problems/move-zeroes/submissions/497455529/
// 输入: nums = [0,1,0,3,12]
// 输出: [1,3,12,0,0]

func moveZeroes(nums []int) {
	if len(nums) == 1 {
		return
	}
	j := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != 0 {
			nums[i], nums[j] = nums[j], nums[i]
			j++
		}
	}
}
