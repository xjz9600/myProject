package leetcode

// https://leetcode.cn/problems/minimum-depth-of-binary-tree/submissions/

func minDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	leftHeight := minDepth(root.Left)
	rightHeight := minDepth(root.Right)
	return minHeight(leftHeight, rightHeight) + 1
}

func minHeight(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
