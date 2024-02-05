package leetcode

// https://leetcode.cn/problems/majority-element/submissions/500482841/

func majorityElement(nums []int) int {
	var sel, count int
	for _, n := range nums {
		if count == 0 {
			sel = n
		} else {
			if n == sel {
				count++
			} else {
				count--
			}
		}
	}
	return sel
}
