package leetcode

//type ListNode struct {
//	Val  int
//	Next *ListNode
//}

// https://leetcode.cn/problems/swap-nodes-in-pairs/

func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	newHead := head.Next
	head.Next = newHead.Next
	newHead.Next = head
	for head.Next != nil && head.Next.Next != nil {
		nextNode := head.Next
		node := nextNode.Next
		head.Next = node
		nextNode.Next = node.Next
		node.Next = nextNode
		head = nextNode
	}
	return newHead
}

func swapPairsV2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	nextNode := swapPairsV2(head.Next.Next)
	prev := head.Next
	head.Next.Next = head
	head.Next = nextNode
	return prev
}
