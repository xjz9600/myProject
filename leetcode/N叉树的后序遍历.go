package leetcode

//https://leetcode.cn/problems/n-ary-tree-postorder-traversal/submissions/498446029/

func postorder(root *Node) []int {
	var result []int
	if root == nil {
		return result
	}
	for _, rc := range root.Children {
		result = append(result, postorder(rc)...)
	}
	result = append(result, root.Val)
	return result
}

func postorderV1(root *Node) []int {
	var result []int
	if root == nil {
		return result
	}
	var stack []*Node
	mapNode := map[*Node]struct{}{}
	stack = append(stack, root)
	for len(stack) > 0 {
		st := stack[0]
		stack = stack[1:]
		if len(st.Children) == 0 {
			result = append(result, st.Val)
			continue
		}
		if _, ok := mapNode[st]; ok {
			result = append(result, st.Val)
			continue
		} else {
			mapNode[st] = struct{}{}
		}
		var nextStack []*Node
		nextStack = append(st.Children, st)
		nextStack = append(nextStack, stack...)
		stack = nextStack
	}
	return result
}
