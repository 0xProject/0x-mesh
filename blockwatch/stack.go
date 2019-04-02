package blockwatch

import (
	"container/list"
	"sync"
)

// Stack allows performing basic stack operations on a stack of MiniHeaders.
type Stack struct {
	limit int
	list  *list.List
	mut   sync.Mutex
}

// NewStack instantiates a new stack with the specified size limit. Once the size limit
// is reached, adding additional blocks will evict the deepest block.
func NewStack(limit int) *Stack {
	return &Stack{
		limit: limit,
		list:  list.New(),
	}
}

// Pop removes and returns the latest block header on the block stack.
func (s *Stack) Pop() *MiniHeader {
	s.mut.Lock()
	defer s.mut.Unlock()
	block := s.list.Front()
	s.list.Remove(block)
	return block.Value.(*MiniHeader)
}

// Push pushes a block header onto the block stack. If the stack limit is reached,
// it will remove the oldest block header and return it.
func (s *Stack) Push(block *MiniHeader) *MiniHeader {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.list.PushFront(block)
	if s.list.Len() > s.limit {
		lastElement := s.list.Back()
		s.list.Remove(lastElement)
		return lastElement.Value.(*MiniHeader)
	}
	return nil
}

// Peek returns the latest block header from the block stack without removing it.
func (s *Stack) Peek() *MiniHeader {
	s.mut.Lock()
	defer s.mut.Unlock()
	block := s.list.Front()
	if block == nil {
		return nil
	}
	return block.Value.(*MiniHeader)
}

// Inspect returns all the block headers currently on the stack. This method should only be
// used for debugging and testing purposes since it is not performant.
func (s *Stack) Inspect() []*MiniHeader {
	s.mut.Lock()
	defer s.mut.Unlock()
	blocks := []*MiniHeader{}
	for e := s.list.Back(); e != nil; e = e.Prev() {
		blocks = append(blocks, e.Value.(*MiniHeader))
	}
	return blocks
}
