package leetcode

import "sort"

//https://leetcode.cn/problems/assign-cookies/description/

//输入: g = [1,2,3], s = [1,1]
//输出: 1
//解释:
//你有三个孩子和两块小饼干，3个孩子的胃口值分别是：1,2,3。
//虽然你有两块小饼干，由于他们的尺寸都是1，你只能让胃口值是1的孩子满足。
//所以你应该输出1。

func findContentChildren(g []int, s []int) int {
	sort.Ints(g)
	sort.Ints(s)
	var num int
	for i, j := 0, 0; i < len(g); i++ {
		for j < len(s) && g[i] > s[j] {
			j++
		}
		if j < len(s) {
			num++
			j++
		} else {
			return num
		}
	}
	return num
}
