package leetcode

// 输入：nums = [-1,-100,3,99], k = 2
// 输出：[3,99,-1,-100]
// 解释:
// 向右轮转 1 步: [99,-1,-100,3]
// 向右轮转 2 步: [3,99,-1,-100]
func rotate(nums []int, k int) {
	var j int
	var nextPrev int
	prev := nums[j]
	if k == 0 {
		return
	}
	var count int
	var start int
	for i := 0; i < len(nums); i++ {
		j = (j + k) % len(nums)
		nextPrev = nums[j]
		nums[j] = prev
		prev = nextPrev
		count++
		if count == len(nums) {
			return
		}
		if j == start {
			j++
			prev = nums[j]
			start = j
		}
	}
}

func rotateV1(nums []int, k int) {
	reverse(nums)
	reverse(nums[:k%len(nums)])
	reverse(nums[k%len(nums):])
}

func reverse(nums []int) {
	for i := 0; i < len(nums)/2; i++ {
		nums[i], nums[len(nums)-1-i] = nums[len(nums)-1-i], nums[i]
	}
}
