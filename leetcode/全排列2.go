package leetcode

// https://leetcode.cn/problems/permutations-ii/submissions/500453556/
// 输入：nums = [1,1,2]
// 输出：
// [[1,1,2],
// [1,2,1],
// [2,1,1]]

func permuteUnique(nums []int) (ans [][]int) {
	var dfs func(re []int, i, k int)
	dfs = func(re []int, n, k int) {
		if n == k {
			prev := make([]int, k)
			copy(prev, re)
			ans = append(ans, prev)
			return
		}
		var visit = map[int]struct{}{}
		for i := n; i < k; i++ {
			if _, ok := visit[nums[i]]; ok {
				continue
			}
			re = append(re, nums[n])
			nums[n], nums[i] = nums[i], nums[n]
			dfs(re, n+1, k)
			nums[i], nums[n] = nums[n], nums[i]
			re = re[:len(re)-1]
			visit[nums[i]] = struct{}{}
		}
	}
	dfs([]int{}, 0, len(nums))
	return
}
