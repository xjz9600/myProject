package leetcode

//https://leetcode.cn/problems/generate-parentheses/

func generateParenthesis(n int) []string {
	return gen(0, 0, n, "")
}

func gen(left, right, n int, s string) []string {
	var result []string
	if left == n && right == n {
		result = append(result, s)
		return result
	}
	if left < n {
		result = append(result, gen(left+1, right, n, s+"(")...)
	}
	if right < left {
		result = append(result, gen(left, right+1, n, s+")")...)
	}
	return result
}
