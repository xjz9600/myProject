package leetcode

// https://leetcode.cn/problems/climbing-stairs/?utm_source=LCUS&utm_medium=ip_redirect&utm_campaign=transfer2china
// 输入：n = 3
// 输出：3
// 解释：有三种方法可以爬到楼顶。
// 1. 1 阶 + 1 阶 + 1 阶
// 2. 1 阶 + 2 阶
// 3. 2 阶 + 1 阶
func climbStairs(n int) int {
	f1 := 1
	f2 := 2
	if n <= 2 {
		return n
	}
	var f3 int
	for i := 3; i <= n; i++ {
		f3 = f1 + f2
		f1 = f2
		f2 = f3
	}
	return f3
}
