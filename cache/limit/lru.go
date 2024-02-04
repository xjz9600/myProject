package limit

type LRUCache[T any] struct {
	data map[string]*Node[T]
	head *Node[T]
	tail *Node[T]
}

type Node[T any] struct {
	preNode  *Node[T]
	nextNode *Node[T]
	value    T
	key      string
}

func ConstructorNode[T any]() LRUCache[T] {
	head := &Node[T]{}
	tail := &Node[T]{}
	head.nextNode = tail
	tail.preNode = head
	return LRUCache[T]{
		data: make(map[string]*Node[T]),
		head: head,
		tail: tail,
	}
}

func (l *LRUCache[T]) RemoveNode(key string) *Node[T] {
	node := l.data[key]
	nextNode := node.nextNode
	node.preNode.nextNode = nextNode
	nextNode.preNode = node.preNode
	node.preNode = nil
	node.nextNode = nil
	delete(l.data, key)
	return node
}

func (l *LRUCache[T]) moveNode(node *Node[T]) {
	nextLNode := node.nextNode
	nextLNode.preNode = node.preNode
	node.preNode.nextNode = nextLNode
	node.nextNode = l.head.nextNode
	l.head.nextNode.preNode = node
	l.head.nextNode = node
	node.preNode = l.head
}

func (l *LRUCache[T]) RemoveTailNode() *Node[T] {
	removeNode := l.tail.preNode
	headNode := removeNode.preNode
	headNode.nextNode = l.tail
	l.tail.preNode = headNode
	removeNode.nextNode = nil
	removeNode.preNode = nil
	delete(l.data, removeNode.key)
	return removeNode
}

func (l *LRUCache[T]) Get(key string) *Node[T] {
	lNode, ok := l.data[key]
	if !ok {
		return nil
	}
	l.moveNode(lNode)
	return lNode
}

func (l *LRUCache[T]) addNewNode(key string, val T) *Node[T] {
	newNode := &Node[T]{
		value: val,
		key:   key,
	}
	l.head.nextNode.preNode = newNode
	newNode.nextNode = l.head.nextNode
	l.head.nextNode = newNode
	newNode.preNode = l.head
	return newNode
}

func (l *LRUCache[T]) Put(key string, value T) {
	lNode, ok := l.data[key]
	if !ok {
		newNode := l.addNewNode(key, value)
		l.data[key] = newNode
		return
	}
	lNode.value = value
	l.moveNode(lNode)
}
