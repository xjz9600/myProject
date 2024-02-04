package leetcode

//https://leetcode.cn/problems/merge-two-sorted-lists/submissions/499574695/

//输入：l1 = [1,2,4], l2 = [1,3,4]
//输出：[1,1,2,3,4,4]

//type ListNode struct {
//	Val  int
//	Next *ListNode
//}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1
	}
	if list1.Val < list2.Val {
		list1.Next = mergeTwoLists(list1.Next, list2)
		return list1
	} else {
		list2.Next = mergeTwoLists(list1, list2.Next)
		return list2
	}
}

func mergeTwoListsV1(list1 *ListNode, list2 *ListNode) *ListNode {
	var merge = &ListNode{Val: 0, Next: nil}
	prev := merge
	for list1 != nil && list2 != nil {
		if list1.Val > list2.Val {
			prev.Next = list2
			list2 = list2.Next
		} else {
			prev.Next = list1
			list1 = list1.Next
		}
		prev = prev.Next
	}
	if list1 != nil {
		prev.Next = list1
	}
	if list2 != nil {
		prev.Next = list2
	}
	return merge.Next
}
