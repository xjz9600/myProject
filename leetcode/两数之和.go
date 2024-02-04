package leetcode

// https://leetcode.cn/problems/two-sum/submissions/497481808/
//输入：nums = [2,7,11,15], target = 9
//输出：[0,1]
//解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。

func twoSum(nums []int, target int) []int {
	numsMap := map[int]int{}
	var result []int
	for k, v := range nums {
		if _, ok := numsMap[target-v]; ok {
			return []int{k, numsMap[target-v]}
		} else {
			numsMap[v] = k
		}
	}
	return result
}
