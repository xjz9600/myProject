package leetcode

//输入：nums = [1,2,3]
//输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]

func permute(nums []int) (ans [][]int) {
	var result []int
	var dfs func(int)
	dfs = func(i int) {
		if len(result) == len(nums) {
			re := make([]int, len(nums))
			copy(re, result)
			ans = append(ans, re)
			return
		}
		result = append(result, nums[i])
		dfs(i + 1)
		result = result[:len(result)-1]
		dfs(i + 1)
	}
	dfs(0)
	return
}
