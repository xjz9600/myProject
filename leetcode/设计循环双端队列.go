package leetcode

// https://leetcode.cn/problems/design-circular-deque/submissions/497995895/?utm_source=LCUS&utm_medium=ip_redirect&utm_campaign=transfer2china
// 输入
// ["MyCircularDeque", "insertLast", "insertLast", "insertFront", "insertFront", "getRear", "isFull", "deleteLast", "insertFront", "getFront"]
// [[3], [1], [2], [3], [4], [], [], [], [4], []]
// 输出
// [null, true, true, true, false, 2, true, true, true, 4]

type MyCircularDeque struct {
	data     []int
	head     int
	tail     int
	capacity int
}

func Constructor(k int) MyCircularDeque {
	return MyCircularDeque{
		data:     make([]int, k+1),
		capacity: k + 1,
	}
}

func (this *MyCircularDeque) InsertFront(value int) bool {
	if this.IsFull() {
		return false
	}
	this.data[this.head] = value
	this.head = (this.head + 1) % this.capacity
	return true
}

func (this *MyCircularDeque) InsertLast(value int) bool {
	if this.IsFull() {
		return false
	}
	this.tail = (this.tail - 1 + this.capacity) % this.capacity
	this.data[this.tail] = value
	return true
}

func (this *MyCircularDeque) DeleteFront() bool {
	if this.IsEmpty() {
		return false
	}
	this.head = (this.head - 1 + this.capacity) % this.capacity
	return true
}

func (this *MyCircularDeque) DeleteLast() bool {
	if this.IsEmpty() {
		return false
	}
	this.tail = (this.tail + 1) % this.capacity
	return true
}

func (this *MyCircularDeque) GetFront() int {
	if this.IsEmpty() {
		return -1
	}
	return this.data[(this.head-1+this.capacity)%this.capacity]
}

func (this *MyCircularDeque) GetRear() int {
	if this.IsEmpty() {
		return -1
	}
	return this.data[this.tail]
}

func (this *MyCircularDeque) IsEmpty() bool {
	return this.head == this.tail
}

func (this *MyCircularDeque) IsFull() bool {
	return (this.head+1)%this.capacity == this.tail
}
