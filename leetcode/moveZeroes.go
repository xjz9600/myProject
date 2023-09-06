package leetcode

func moveZeroes(nums []int) {
	var j int
	for i, n := range nums {
		if n != 0 {
			nums[i], nums[j] = nums[j], nums[i]
			j++
		}
	}
}
