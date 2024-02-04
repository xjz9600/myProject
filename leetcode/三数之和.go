package leetcode

import "sort"

// https://leetcode.cn/problems/3sum/solutions/284681/san-shu-zhi-he-by-leetcode-solution/
//输入：nums = [-1,0,1,2,-1,-4]
//输出：[[-1,-1,2],[-1,0,1]]
//解释：
//nums[0] + nums[1] + nums[2] = (-1) + 0 + 1 = 0 。
//nums[1] + nums[2] + nums[4] = 0 + 1 + (-1) = 0 。
//nums[0] + nums[3] + nums[4] = (-1) + 2 + (-1) = 0 。
//不同的三元组是 [-1,0,1] 和 [-1,-1,2] 。
//注意，输出的顺序和三元组的顺序并不重要。

// 先排序
//nums = [-2,0,0,2,2]
//nums = [-4,-1,-1,0,1,2]

func threeSum(nums []int) [][]int {
	result := [][]int{}
	sort.Ints(nums)
	for i := 0; i < len(nums)-2; i++ {
		first := nums[i]
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		j := i + 1
		k := len(nums) - 1
		for k > j {
			if j > i+1 && nums[j] == nums[j-1] {
				j++
				continue
			}
			if first+nums[j]+nums[k] < 0 {
				j++
			} else if first+nums[j]+nums[k] > 0 {
				k--
			} else {
				re := []int{first, nums[j], nums[k]}
				result = append(result, re)
				j++
				k--
			}
		}
	}
	return result
}
