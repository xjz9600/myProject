package leetcode

// https://leetcode.cn/problems/subsets/submissions/500478254/
// 输入：nums = [1,2,3]
// 输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
func subsets(nums []int) (ans [][]int) {
	var dfs func(re []int, i int)
	dfs = func(re []int, i int) {
		if i > len(nums) {
			return
		}
		result := make([]int, len(re))
		copy(result, re)
		ans = append(ans, result)
		for j := i; j < len(nums); j++ {
			re = append(re, nums[j])
			dfs(re, j+1)
			re = re[:len(re)-1]
		}
	}
	dfs([]int{}, 0)
	return
}
