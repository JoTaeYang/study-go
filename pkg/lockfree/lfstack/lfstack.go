package lfstack

import "sync/atomic"

type Node[T any] struct {
	value T
	next  *Node[T]
}

type TopNode[T any] struct {
	next   *Node[T]
	unique int64
}

type Stack[T any] struct {
	top   atomic.Pointer[TopNode[T]] //nil이 되지 않는다.
	count atomic.Int32
}

func NewStack[T any]() *Stack[T] {
	s := Stack[T]{}
	s.top.Store(&TopNode[T]{
		next:   nil,
		unique: 0,
	})
	return &s
}

func (s *Stack[T]) Push(val T) {
	newNode := &Node[T]{value: val}

	var tmpTop *TopNode[T]

	tmpUnique := atomic.AddInt64(&s.top.Load().unique, 1)

	//var newTop TopNode[T]
	for {
		tmpTop = s.top.Load()

		newTop := &TopNode[T]{
			next:   newNode,
			unique: tmpUnique,
		}

		newNode.next = tmpTop.next
		if s.top.CompareAndSwap(tmpTop, newTop) {
			break
		}
	}
	s.count.Add(1)
}

func (s *Stack[T]) GetCount() int32 {
	return s.count.Load()
}

func (s *Stack[T]) CountDecrement() bool {
	if s.count.Add(-1) < 0 {
		s.count.Add(1)
		return false
	}
	return true
}

func (s *Stack[T]) Pop() (T, bool) {
	var empty T

	tmpUnique := atomic.AddInt64(&s.top.Load().unique, 1)
	if !s.CountDecrement() {
		return empty, false
	}
	var oldTop *TopNode[T]
	//var newTop TopNode[T]
	for {
		oldTop = s.top.Load()
		if oldTop.next == nil {
			return empty, false
		}

		newTop := &TopNode[T]{
			next:   oldTop.next.next,
			unique: tmpUnique,
		}

		if s.top.CompareAndSwap(oldTop, newTop) {
			return oldTop.next.value, true
		}
	}
}
