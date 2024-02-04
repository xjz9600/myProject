package leetcode

// https://leetcode.cn/problems/UHnkqh/submissions/497497775/

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	nextNode := head.Next
	head.Next = nil
	for nextNode != nil {
		node := nextNode.Next
		nextNode.Next = head
		head = nextNode
		nextNode = node
	}
	return head
}

func reverseListV2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	newHeader := reverseList(head.Next)
	head.Next.Next = head
	head.Next = nil
	return newHeader
}
