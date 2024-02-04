package leetcode

// https://leetcode.cn/problems/binary-tree-preorder-traversal/submissions/498181352/

func preorderTraversal(root *TreeNode) []int {
	var result []int
	if root == nil {
		return result
	}
	result = append(result, root.Val)
	leftResult := preorderTraversal(root.Left)
	result = append(result, leftResult...)
	rightResult := preorderTraversal(root.Right)
	result = append(result, rightResult...)
	return result
}

// [5,3,null,4,2,null,1]

func preorderTraversalV1(root *TreeNode) []int {
	var result []int
	if root == nil {
		return result
	}
	var stack []*TreeNode
	for root != nil || len(stack) > 0 {
		for root != nil {
			result = append(result, root.Val)
			stack = append(stack, root)
			root = root.Left
		}
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		root = st.Right
	}
	return result
}
