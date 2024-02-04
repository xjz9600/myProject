package leetcode

import "math"

//https://leetcode.cn/problems/validate-binary-search-tree/description/
//给你一个二叉树的根节点 root ，判断其是否是一个有效的二叉搜索树。
//
//有效 二叉搜索树定义如下：
//
//节点的左子树只包含 小于 当前节点的数。
//节点的右子树只包含 大于 当前节点的数。
//所有左子树和右子树自身必须也是二叉搜索树。

func isValidBST(root *TreeNode) bool {
	min := math.MinInt
	max := math.MaxInt
	return isValidBSTOK(root, min, max)
}

func isValidBSTOK(node *TreeNode, min, max int) bool {
	if node == nil {
		return true
	}
	if node.Val < min || node.Val > max {
		return false
	}
	leftBST := isValidBSTOK(node.Left, min, minData(max, node.Val)-1)
	rightBST := isValidBSTOK(node.Right, maxData(min, node.Val)+1, max)
	return leftBST && rightBST
}

func minData(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func maxData(a, b int) int {
	if a <= b {
		return b
	}
	return a
}
