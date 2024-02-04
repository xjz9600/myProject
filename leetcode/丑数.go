package leetcode

//https://leetcode.cn/problems/chou-shu-lcof/submissions/499074306/

//给你一个整数 n ，请你找出并返回第 n 个 丑数 。

//说明：丑数是只包含质因数 2、3、5 的正整数；1 是丑数。

//输入: n = 10
//输出: 12
//解释: 1, 2, 3, 4, 5, 6, 8, 9, 10, 12 是前 10 个丑数。

func nthUglyNumber(n int) int {
	var heap = []int{1}
	var mapHeap = map[int]struct{}{}
	moveTail := func() {
		i := len(heap) - 1
		j := (i - 1) / 2
		for j >= 0 && heap[i] < heap[j] {
			heap[i], heap[j] = heap[j], heap[i]
			i = j
			j = (i - 1) / 2
		}
	}
	moveHeader := func() {
		i := 0
		lenHeap := len(heap) - 1
		for 2*i+1 <= lenHeap || 2*i+2 <= lenHeap {
			ni := 2*i + 1
			if 2*i+2 <= lenHeap {
				if heap[2*i+2] < heap[2*i+1] {
					ni = 2*i + 2
				}
			}
			if heap[ni] >= heap[i] {
				break
			}
			heap[ni], heap[i] = heap[i], heap[ni]
			i = ni
		}
	}
	push := func(k int) {
		var arr = []int{2, 3, 5}
		for _, a := range arr {
			if _, ok := mapHeap[k*a]; ok {
				continue
			}
			heap = append(heap, k*a)
			mapHeap[k*a] = struct{}{}
			moveTail()
		}
	}
	pop := func() int {
		val := heap[0]
		heap[0] = heap[len(heap)-1]
		heap = heap[:len(heap)-1]
		moveHeader()
		return val
	}
	var startVal int
	for i := 1; i < n; i++ {
		startVal = pop()
		push(startVal)
	}
	return heap[0]
}
