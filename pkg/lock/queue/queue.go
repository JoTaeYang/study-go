package queue

import (
	"sync"
)

type Node[T any] struct {
	value T
	next  *Node[T]
}

type Queue[T any] struct {
	head  *Node[T]
	tail  *Node[T]
	count int32
	qlock sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{

		qlock: sync.Mutex{},
	}

}

func (q *Queue[T]) Enqueue(val T) {
	q.qlock.Lock()
	defer q.qlock.Unlock()

	newNode := &Node[T]{value: val, next: nil}

	if q.tail == nil {
		q.tail = newNode
		q.head = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}

	q.count++
}

func (q *Queue[T]) Dequeue() T {
	q.qlock.Lock()
	defer q.qlock.Unlock()

	var empty T
	if q.head == nil {
		return empty
	}
	q.count--
	if q.count < 0 {
		q.count++
		return empty
	}

	outValue := q.head.value
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}

	return outValue
}

func (q *Queue[T]) Peek() T {
	q.qlock.Lock()
	defer q.qlock.Unlock()

	var empty T
	if q.head == nil {
		return empty
	}

	outValue := q.head.value

	return outValue
}

func (q *Queue[T]) GetCount() int32 {
	return q.count
}
