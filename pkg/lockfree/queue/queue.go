package queue

import (
	"sync/atomic"
	"unsafe"
)

type Node[T any] struct {
	value T
	next  *Node[T]
}

type TopNode[T any] struct {
	node   *Node[T]
	unique int64
}

type Queue[T any] struct {
	head   atomic.Pointer[TopNode[T]]
	tail   atomic.Pointer[TopNode[T]]
	unique atomic.Int64
	count  atomic.Int32
}

func NewQueue[T any]() *Queue[T] {
	q := Queue[T]{}
	dummyNode := Node[T]{
		next: nil,
	}

	head := TopNode[T]{
		node:   &dummyNode,
		unique: 0,
	}

	tail := TopNode[T]{
		node:   head.node,
		unique: 0,
	}

	q.count.Store(0)
	q.unique.Store(0)
	q.head.Store(&head)
	q.tail.Store(&tail)

	return &q
}

func (q *Queue[T]) Enqueue(val T) {
	newNode := &Node[T]{value: val}
	tmpUnique := q.unique.Add(1)
	var tmpTail *TopNode[T]

	for {
		tmpTail = q.tail.Load()

		if tmpTail.node.next == nil {
			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(tmpTail.node.next)),
				nil,
				unsafe.Pointer(newNode)) {
				if q.tail.CompareAndSwap(tmpTail, &TopNode[T]{
					node:   tmpTail.node.next,
					unique: tmpUnique,
				}) {
					break
				}
			}
		} else {
			q.tail.CompareAndSwap(tmpTail, &TopNode[T]{
				node:   tmpTail.node.next,
				unique: tmpUnique,
			})
		}
	}
	q.count.Add(1)
}

func (q *Queue[T]) Dequeue(val *T) {
	var tmpTail *TopNode[T]
	var tmpHead *TopNode[T]

	if q.count.Add(-1) < 0 {
		q.count.Add(1)
		return
	}

	tmpUnique := q.unique.Add(1)
	for {
		tmpTail = q.tail.Load()
		tmpHead = q.head.Load()
		if tmpTail.node.next == nil {
			if tmpHead.node.next != nil {
				val = &tmpHead.node.value

				if q.head.CompareAndSwap(tmpHead, &TopNode[T]{
					node:   tmpHead.node.next,
					unique: tmpUnique,
				}) {

				}
			}
		} else {
			q.tail.CompareAndSwap(tmpTail, &TopNode[T]{
				node:   tmpTail.node.next,
				unique: tmpUnique,
			})
		}
	}
}
