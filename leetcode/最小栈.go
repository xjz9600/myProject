package leetcode

// https://leetcode.cn/problems/min-stack/submissions/497665253/
//["MinStack","push","push","push","getMin","pop","top","getMin"]
//[[],[-2],[0],[-3],[],[],[],[]]
//
//输出：
//[null,null,null,null,-3,null,0,-2]

type MinStack struct {
	stack    []int
	minStack []int
}

func ConstructorV1() MinStack {
	return MinStack{
		stack:    []int{},
		minStack: []int{},
	}
}

func minVal(a, b int) int {
	if a >= b {
		return b
	}
	return a
}

func (this *MinStack) Push(val int) {
	this.stack = append(this.stack, val)
	if len(this.minStack) == 0 {
		this.minStack = append(this.minStack, val)
	} else {
		this.minStack = append(this.minStack, minVal(this.minStack[len(this.minStack)-1], val))
	}
}

func (this *MinStack) Pop() {
	this.stack = this.stack[:len(this.stack)-1]
	this.minStack = this.minStack[:len(this.minStack)-1]
}

func (this *MinStack) Top() int {
	return this.stack[len(this.stack)-1]
}

func (this *MinStack) GetMin() int {
	return this.minStack[len(this.minStack)-1]
}
