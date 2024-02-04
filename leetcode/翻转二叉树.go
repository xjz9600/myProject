package leetcode

//https://leetcode.cn/problems/invert-binary-tree/submissions/499658475/

//type TreeNode struct {
//	Val   int
//	Left  *TreeNode
//	Right *TreeNode
//}

func invertTree(root *TreeNode) *TreeNode {
	if root == nil || (root.Left == nil && root.Right == nil) {
		return root
	}
	oldLeft := root.Left
	root.Left = invertTree(root.Right)
	root.Right = invertTree(oldLeft)
	return root
}
