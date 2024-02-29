package leetcode

//https://leetcode.cn/problems/search-in-rotated-sorted-array/description/

//输入：nums = [4,5,6,7,0,1,2], target = 0
//输出：4

func search(nums []int, target int) int {
	if len(nums) == 1 {
		if nums[0] == target {
			return 0
		} else {
			return -1
		}
	}
	l := 0
	r := len(nums) - 1
	for l <= r {
		mid := (l + r) / 2
		if nums[mid] == target {
			return mid
		}
		if nums[l] > nums[mid] {
			if nums[mid] < target && target <= nums[r] {
				l = mid + 1
			} else {
				r = mid - 1
			}
		} else {
			if nums[l] <= target && nums[mid] > target {
				r = mid - 1
			} else {
				l = mid + 1
			}
		}
	}
	return -1
}
