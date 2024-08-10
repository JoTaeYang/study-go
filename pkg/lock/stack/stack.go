package lock

import (
	"sync"
)

type Node[T any] struct {
	value T
	next  *Node[T]
}

type Stack[T any] struct {
	top   *Node[T]
	count int32
	mu    sync.Mutex
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newNode := &Node[T]{value: val}
	newNode.next = s.top
	s.top = newNode
	s.count++
}

func (s *Stack[T]) GetCount() int32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count
}

func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var empty T

	if s.top == nil {
		return empty, false
	}

	value := s.top.value
	s.top = s.top.next
	s.count--

	return value, true
}
