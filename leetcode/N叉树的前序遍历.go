package leetcode

// https://leetcode.cn/problems/n-ary-tree-preorder-traversal/submissions/498443007/

type Node struct {
	Val      int
	Children []*Node
}

func preorder(root *Node) []int {
	var result []int
	if root == nil {
		return result
	}
	result = append(result, root.Val)
	for _, rc := range root.Children {
		result = append(result, preorder(rc)...)
	}
	return result
}

func preorderV1(root *Node) []int {
	var result []int
	if root == nil {
		return result
	}
	var stack []*Node
	stack = append(stack, root)
	for len(stack) > 0 {
		st := stack[0]
		stack = stack[1:]
		result = append(result, st.Val)
		newStack := st.Children
		newStack = append(newStack, stack...)
		stack = newStack
	}
	return result
}
