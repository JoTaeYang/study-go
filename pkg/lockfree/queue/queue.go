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
	dummyNode := &Node[T]{
		next: nil,
	}

	// head :=

	// tail :=

	q.count.Store(0)
	q.unique.Store(0)
	q.head.Store(&TopNode[T]{
		node:   dummyNode,
		unique: 0,
	})
	q.tail.Store(&TopNode[T]{
		node:   q.head.Load().node,
		unique: 0,
	})

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
				(*unsafe.Pointer)(unsafe.Pointer(&tmpTail.node.next)),
				nil,
				unsafe.Pointer(&newNode)) {

				newTopNode := &TopNode[T]{
					node:   tmpTail.node.next,
					unique: tmpUnique,
				}

				if q.tail.CompareAndSwap(tmpTail, newTopNode) {
					break
				}
			}
		} else {
			newTopN := &TopNode[T]{
				node:   tmpTail.node.next,
				unique: tmpUnique,
			}
			q.tail.CompareAndSwap(tmpTail, newTopN)
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
				이거 값이 좀 잘못들어가는 거 같은데 다시 파악해보기
				val = &tmpHead.node.next.value

				if q.head.CompareAndSwap(tmpHead, &TopNode[T]{
					node:   tmpHead.node.next,
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
}
