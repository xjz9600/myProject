package leetcode

//https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-inorder-traversal/submissions/500049054/
//输入: preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
//输出: [3,9,20,null,null,15,7]

func buildTree(preorder []int, inorder []int) *TreeNode {
	node := &TreeNode{Val: preorder[0]}
	if len(preorder) == 1 {
		return node
	}
	middle := preorder[0]
	var left int
	for j, in := range inorder {
		if in == middle {
			left = j
			break
		}
	}
	if left+1 > 1 {
		node.Left = buildTree(preorder[1:left+1], inorder[:left])
	}
	if left+1 < len(preorder) {
		node.Right = buildTree(preorder[left+1:], inorder[left+1:])
	}
	return node
}
