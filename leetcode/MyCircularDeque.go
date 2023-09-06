package leetcode

type MyCircularDeque struct {
	data  []int
	front int
	last  int
	count int
}

func Constructor(k int) MyCircularDeque {
	return MyCircularDeque{
		data: make([]int, k),
	}
}

func (this *MyCircularDeque) InsertFront(value int) bool {
	if this.IsFull() {
		return false
	}
	this.count++
	this.front = (len(this.data) + this.front - 1) % len(this.data)
	this.data[this.front] = value
	return true
}

func (this *MyCircularDeque) InsertLast(value int) bool {
	if this.IsFull() {
		return false
	}
	this.count++
	this.data[this.last] = value
	this.last = (this.last + 1) % len(this.data)
	return true
}

func (this *MyCircularDeque) DeleteFront() bool {
	if this.IsEmpty() {
		return false
	}
	this.count--
	this.front = (this.front + 1) % len(this.data)
	return true
}

func (this *MyCircularDeque) DeleteLast() bool {
	if this.IsEmpty() {
		return false
	}
	this.count--
	this.last = (this.last - 1 + len(this.data)) % len(this.data)
	return true
}

func (this *MyCircularDeque) GetFront() int {
	if this.IsEmpty() {
		return -1
	}
	return this.data[this.front]
}

func (this *MyCircularDeque) GetRear() int {
	if this.IsEmpty() {
		return -1
	}
	return this.data[(this.last-1+len(this.data))%len(this.data)]
}

func (this *MyCircularDeque) IsEmpty() bool {
	return this.count == 0
}

func (this *MyCircularDeque) IsFull() bool {
	return this.count == cap(this.data)
}

// 标准答案
//type MyCircularDeque struct {
//	data  []int
//	front int
//	last  int
//}
//
//func Constructor(k int) MyCircularDeque {
//	return MyCircularDeque{
//		data: make([]int, k+1),
//	}
//}
//
//func (this *MyCircularDeque) InsertFront(value int) bool {
//	if this.IsFull() {
//		return false
//	}
//	this.front = (len(this.data) + this.front - 1) % len(this.data)
//	this.data[this.front] = value
//	return true
//}
//
//func (this *MyCircularDeque) InsertLast(value int) bool {
//	if this.IsFull() {
//		return false
//	}
//	this.data[this.last] = value
//	this.last = (this.last + 1) % len(this.data)
//	return true
//}
//
//func (this *MyCircularDeque) DeleteFront() bool {
//	if this.IsEmpty() {
//		return false
//	}
//	this.front = (this.front + 1) % len(this.data)
//	return true
//}
//
//func (this *MyCircularDeque) DeleteLast() bool {
//	if this.IsEmpty() {
//		return false
//	}
//	this.last = (this.last - 1 + len(this.data)) % len(this.data)
//	return true
//}
//
//func (this *MyCircularDeque) GetFront() int {
//	if this.IsEmpty() {
//		return -1
//	}
//	return this.data[this.front]
//}
//
//func (this *MyCircularDeque) GetRear() int {
//	if this.IsEmpty() {
//		return -1
//	}
//	return this.data[(this.last-1+len(this.data))%len(this.data)]
//}
//
//func (this *MyCircularDeque) IsEmpty() bool {
//	return this.front == this.last
//}
//
//func (this *MyCircularDeque) IsFull() bool {
//	return (this.last+1)%len(this.data) == this.front
//}
