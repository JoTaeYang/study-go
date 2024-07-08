package lfstack

import "sync/atomic"

type Node struct {
	value int
	next  *Node
}

type Stack struct {
	top atomic.Pointer[Node]
}

func (s *Stack) Push(val int) {
	newNode := &Node{value: val}
	for {
		oldTop := s.top.Load()
		newNode.next = oldTop
		if s.top.CompareAndSwap(oldTop, newNode) {
			break
		}
	}
}

func (s *Stack) Pop() (int, bool) {
	for {
		oldTop := s.top.Load()
		if oldTop == nil {
			return 0, false
		}
		newTop := oldTop.next
		if s.top.CompareAndSwap(oldTop, newTop) {
			return oldTop.value, true
		}
	}
}
