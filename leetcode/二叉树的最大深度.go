package leetcode

// https://leetcode.cn/problems/maximum-depth-of-binary-tree/submissions/499837323/
//给定一个二叉树 root ，返回其最大深度。

//二叉树的 最大深度 是指从根节点到最远叶子节点的最长路径上的节点数。

func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	left := maxDepth(root.Left)
	right := maxDepth(root.Right)
	return MaxHeight(left, right) + 1
}

func MaxHeight(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
