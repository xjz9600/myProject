package leetcode

// https://leetcode.cn/problems/linked-list-cycle/description/
func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	slow := head
	fast := head.Next
	for slow != fast {
		slow = slow.Next
		if fast.Next == nil || fast.Next.Next == nil {
			return false
		}
		fast = fast.Next.Next
	}
	return true
}
