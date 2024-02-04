package leetcode

// 注：前、中、后序遍历针对的是父节点
// https://leetcode.cn/problems/binary-tree-inorder-traversal/submissions/498182117/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorderTraversal(root *TreeNode) []int {
	var result []int
	if root == nil {
		return result
	}
	leftResult := inorderTraversal(root.Left)
	result = append(result, leftResult...)
	result = append(result, root.Val)
	rightResult := inorderTraversal(root.Right)
	result = append(result, rightResult...)
	return result
}

func inorderTraversalV1(root *TreeNode) []int {
	var result []int
	if root == nil {
		return result
	}
	var stack []*TreeNode
	for root != nil || len(stack) > 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, st.Val)
		root = st.Right
	}
	return result
}
