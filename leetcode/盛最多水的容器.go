package leetcode

// https://leetcode.cn/problems/container-with-most-water/description/

//输入：[1,8,6,2,5,4,8,3,7]
//输出：49

// 题解：因为容器的容量是按照最低的柱子来保证的，所以使用双指针移动较小的柱子来保证

func maxArea(height []int) int {
	i := 0
	j := len(height) - 1
	maxLeft := height[i]
	maxRight := height[j]
	maxSize := (j - i) * min(maxLeft, maxRight)
	for j > i {
		if maxLeft > maxRight {
			j--
			if height[j] > maxRight {
				maxRight = height[j]
				maxSize = max(maxSize, (j-i)*min(maxLeft, maxRight))
			}
		} else {
			i++
			if height[i] > maxLeft {
				maxLeft = height[i]
				maxSize = max(maxSize, (j-i)*min(maxLeft, maxRight))
			}
		}
	}
	return maxSize
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
