package leetcode

//https://leetcode.cn/problems/jump-game/description/

//输入：nums = [2,3,1,1,4]
//输出：true
//解释：可以先跳 1 步，从下标 0 到达下标 1, 然后再从下标 1 跳 3 步到达最后一个下标。

func canJump(nums []int) bool {
	var mostRight int
	for i := 0; i < len(nums); i++ {
		if i <= mostRight {
			mostRight = max(mostRight, nums[i]+i)
			if mostRight >= len(nums)-1 {
				return true
			}
		} else {
			break
		}
	}
	return false
}

func canJump2(nums []int) bool {
	var mostRight int
	var end int
	for i := 0; i < len(nums)-1; i++ {
		mostRight = max(mostRight, nums[i]+i)
		if end == i {
			if mostRight == end {
				return false
			}
			end = mostRight
		}

	}
	return true
}
