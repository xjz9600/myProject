package leetcode

//https://leetcode.cn/problems/permutations/submissions/500412571/
//输入：nums = [1,2,3]
//输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]

func permute(nums []int) (ans [][]int) {
	var re []int
	var genNums func(re []int, dept, k int, nums []int)
	genNums = func(re []int, dept, k int, nums []int) {
		if dept == k {
			result := make([]int, k)
			copy(result, re)
			ans = append(ans, result)
			return
		}
		for i := dept; i < k; i++ {
			re = append(re, nums[i])
			nums[i], nums[dept] = nums[dept], nums[i]
			genNums(re, dept+1, k, nums)
			nums[dept], nums[i] = nums[i], nums[dept]
			re = re[:len(re)-1]
		}
	}
	genNums(re, 0, len(nums), nums)
	return ans
}
