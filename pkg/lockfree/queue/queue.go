package queue

import "sync/atomic"

type Node[T any] struct {
	value T
	next  *Node[T]
}

type TopNode[T any] struct {
	next   *Node[T]
	unique int64
}

type Queue[T any] struct {
	head   atomic.Pointer[TopNode[T]]
	tail   atomic.Pointer[TopNode[T]]
	unique atomic.Int64
	count  atomic.Int32
}

func (q *Queue[T]) Enqueue(val T) {
	newNode := &Node[T]{value: val}
	tmpUnique := q.unique.Add(1)
	var tmpTail *TopNode[T]
	for {
		tmpTail = q.tail.Load()

		if tmpTail.next.next == nil {

		}
	}
}
