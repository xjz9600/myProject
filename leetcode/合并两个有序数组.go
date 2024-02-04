package leetcode

// https://leetcode.cn/problems/merge-sorted-array/submissions/499597497/
// 输入：nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
// 输出：[1,2,2,3,5,6]
// 解释：需要合并 [1,2,3] 和 [2,5,6] 。
// 合并结果是 [1,2,2,3,5,6] ，其中斜体加粗标注的为 nums1 中的元素

func merge(nums1 []int, m int, nums2 []int, n int) {
	var sortNums = make([]int, 0, m+n)
	p1, p2 := 0, 0
	for {
		if p1 == m {
			sortNums = append(sortNums, nums2[p2:]...)
			break
		}
		if p2 == n {
			sortNums = append(sortNums, nums1[p1:]...)
			break
		}
		if nums1[p1] < nums2[p2] {
			sortNums = append(sortNums, nums1[p1])
			p1++
		} else {
			sortNums = append(sortNums, nums2[p2])
			p2++
		}
	}
	copy(nums1, sortNums)
}

func mergeV1(nums1 []int, m int, nums2 []int, n int) {
	p1 := m - 1
	p2 := n - 1
	j := len(nums1)
	for p1 >= -1 && p2 >= -1 {
		j--
		if p2 == -1 {
			return
		}
		if p1 == -1 {
			nums1[j] = nums2[p2]
			p2--
		}
		if nums1[p1] < nums2[p2] {
			nums1[j] = nums2[p2]
			p2--
		} else {
			nums1[j] = nums1[p1]
			p1--
		}
	}
}
