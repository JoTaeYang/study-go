package lfstack

import "sync/atomic"

type Node[T any] struct {
	value T
	next  *Node[T]
}

type Stack[T any] struct {
	top atomic.Pointer[Node[T]]
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
			return oldTop.value, true
		}
	}
}
