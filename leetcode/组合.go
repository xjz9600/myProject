package leetcode

// https://leetcode.cn/problems/combinations/submissions/500081913/

// 输入：n = 4, k = 2
// 输出：
// [
// [2,4],
// [3,4],
// [2,3],
// [1,2],
// [1,3],
// [1,4],
// ]
func combine(n int, k int) (ans [][]int) {
	var result []int
	var dfs func(cur int)
	dfs = func(cur int) {
		if len(result)+(n-cur)+1 < k {
			return
		}
		if len(result) == k {
			re := make([]int, k)
			copy(re, result)
			ans = append(ans, re)
			return
		}
		result = append(result, cur)
		dfs(cur + 1)
		result = result[:len(result)-1]
		dfs(cur + 1)
	}
	dfs(1)
	return ans
}
