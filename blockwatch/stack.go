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
func (bs *Stack) Pop() *MiniHeader {
	bs.mut.Lock()
	defer bs.mut.Unlock()
	block := bs.list.Front()
	bs.list.Remove(block)
	return block.Value.(*MiniHeader)
}

// Push pushes a block header onto the block stack.
func (bs *Stack) Push(block *MiniHeader) {
	bs.mut.Lock()
	defer bs.mut.Unlock()
	bs.list.PushFront(block)
	if bs.list.Len() > bs.limit {
		lastElement := bs.list.Back()
		bs.list.Remove(lastElement)
	}
}

// Peek returns the latest block header from the block stack without removing it.
func (bs *Stack) Peek() *MiniHeader {
	bs.mut.Lock()
	defer bs.mut.Unlock()
	block := bs.list.Front()
	if block == nil {
		return nil
	}
	return block.Value.(*MiniHeader)
}

// Inspect returns all the block headers currently on the stack. This method should only be
// used for debugging and testing purposes since it is not performant.
func (bs *Stack) Inspect() []*MiniHeader {
	bs.mut.Lock()
	defer bs.mut.Unlock()
	blocks := []*MiniHeader{}
	for e := bs.list.Back(); e != nil; e = e.Prev() {
		blocks = append(blocks, e.Value.(*MiniHeader))
	}
	return blocks
}
