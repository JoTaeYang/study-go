package lfstack

import "sync/atomic"

type Node[T any] struct {
	value T
	next  *Node[T]
}

type Stack[T any] struct {
	top   atomic.Pointer[Node[T]]
	count atomic.Int32
}

func (s *Stack[T]) Push(val T) {
	newNode := &Node[T]{value: val}
	for {
		oldTop := s.top.Load()
		newNode.next = oldTop
		if s.top.CompareAndSwap(oldTop, newNode) {
			break
		}
	}
	s.count.Add(1)
}

func (s *Stack[T]) GetCount() int32 {
	return s.count.Load()
}

func (s *Stack[T]) Pop() (T, bool) {
	var empty T
	for {
		oldTop := s.top.Load()
		if oldTop == nil {
			return empty, false
		}
		newTop := oldTop.next
		if s.top.CompareAndSwap(oldTop, newTop) {
			s.count.Add(-1)
			return oldTop.value, true
		}
	}

}
