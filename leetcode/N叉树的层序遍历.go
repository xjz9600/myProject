package leetcode

//https://leetcode.cn/problems/n-ary-tree-level-order-traversal/description/

func levelOrder(root *Node) [][]int {
	var result [][]int
	if root == nil {
		return result
	}
	var stack []*Node
	stack = append(stack, root)
	for len(stack) > 0 {
		var re []int
		var nextStack []*Node
		for _, st := range stack {
			re = append(re, st.Val)
			nextStack = append(nextStack, st.Children...)
		}
		result = append(result, re)
		stack = nextStack
	}
	return result
}
