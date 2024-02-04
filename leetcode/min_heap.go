package leetcode

// 使用小顶堆实现
func findKthLargest(nums []int, k int) int {
	heap := make([]int, 0)
	heap = append(heap, nums[0])
	for i := 1; i < k; i++ {
		heap = genMinHeap(heap, nums[i])
	}
	for i := k; i < len(nums); i++ {
		if nums[i] <= heap[0] {
			continue
		} else {
			heap[0] = nums[i]
			heap = removeHead(heap)
		}
	}
	return heap[0]
}

func removeHead(num []int) []int {
	i := 0
	k := len(num)
	for 2*i+1 < k {
		minNum := num[2*i+1]
		nextI := 2*i + 1
		if 2*i+2 < k {
			if num[2*i+1] > num[2*i+2] {
				minNum = num[2*i+2]
				nextI = 2*i + 2
			}
		}
		if num[i] > minNum {
			num[nextI], num[i] = num[i], num[nextI]
			i = nextI
		} else {
			return num
		}
	}
	return num
}

func genMinHeap(num []int, n int) []int {
	num = append(num, n)
	i := len(num) - 1
	for i > 0 {
		if i%2 == 0 {
			if num[i] < num[i/2-1] {
				num[i], num[i/2-1] = num[i/2-1], num[i]
				i = i/2 - 1
				continue
			}
		} else {
			if num[i] < num[i/2] {
				num[i], num[i/2] = num[i/2], num[i]
				i = i / 2
				continue
			}
		}
		break
	}
	return num
}

// 给定整数数组 nums 和整数 k，请返回数组中第 k 个最大的元素。
//
// 请注意，你需要找的是数组排序后的第 k 个最大的元素，而不是第 k 个不同的元素。
//
// 你必须设计并实现时间复杂度为 O(n) 的算法解决此问题。
//
// 示例 1:
//
// 输入: [3,2,1,5,6,4], k = 2
// 输出: 5
// 示例 2:
//
// 输入: [3,2,3,1,2,4,5,5,6], k = 4
// 输出: 4
//
// 提示：
//
// 1 <= k <= nums.length <= 105
// -104 <= nums[i] <= 104
// 使用快排+快慢指针实现
func findKthLargestV2(nums []int, k int) int {
	return GetMaxK(0, len(nums)-1, k, nums)
}

func GetMaxK(l, r, k int, nums []int) int {
	if l == r {
		return nums[l]
	}
	val := nums[l]
	i := l
	s := i - 1
	b := r + 1
	for i < b {
		if nums[i] < val {
			s++
			nums[i], nums[s] = nums[s], nums[i]
			i++
		} else if nums[i] == val {
			i++
		} else {
			b--
			nums[i], nums[b] = nums[b], nums[i]
		}
	}
	if r-b+1 >= k {
		return GetMaxK(b, r, k, nums)
	}
	if r-s < k {
		return GetMaxK(l, s, k-(r-s), nums)
	}
	return val
}
