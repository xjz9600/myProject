package leetcode

type LRUCache struct {
	data     map[int]*LRUNode
	head     *LRUNode
	tail     *LRUNode
	capacity int
}

type LRUNode struct {
	preNode  *LRUNode
	nextNode *LRUNode
	value    int
	key      int
}

func LruCacheConstructor(capacity int) LRUCache {
	head := &LRUNode{value: -1}
	tail := &LRUNode{value: -1}
	head.nextNode = tail
	tail.preNode = head
	return LRUCache{
		data:     make(map[int]*LRUNode, capacity),
		head:     head,
		tail:     tail,
		capacity: capacity,
	}
}

func (this *LRUCache) moveNode(node *LRUNode) {
	nextLNode := node.nextNode
	nextLNode.preNode = node.preNode
	node.preNode.nextNode = nextLNode
	node.nextNode = this.head.nextNode
	this.head.nextNode.preNode = node
	this.head.nextNode = node
	node.preNode = this.head
}

func (this *LRUCache) removeNode() {
	removeNode := this.tail.preNode
	headNode := removeNode.preNode
	headNode.nextNode = this.tail
	this.tail.preNode = headNode
	removeNode.nextNode = nil
	removeNode.preNode = nil
	delete(this.data, removeNode.key)
}

func (this *LRUCache) Get(key int) int {
	lNode, ok := this.data[key]
	if !ok {
		return -1
	}
	this.moveNode(lNode)
	return lNode.value
}

func (this *LRUCache) addNewNode(key, val int) *LRUNode {
	newNode := &LRUNode{
		value: val,
		key:   key,
	}
	this.head.nextNode.preNode = newNode
	newNode.nextNode = this.head.nextNode
	this.head.nextNode = newNode
	newNode.preNode = this.head
	return newNode
}

func (this *LRUCache) Put(key int, value int) {
	lNode, ok := this.data[key]
	if !ok {
		newNode := this.addNewNode(key, value)
		if len(this.data) == this.capacity {
			this.removeNode()
		}
		this.data[key] = newNode
		return
	}
	lNode.value = value
	this.moveNode(lNode)
}
