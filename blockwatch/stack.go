package blockwatch

import (
	"container/list"
	"sync"
)

// Stack allows performing basic stack operations on a stack of MiniHeaders.
type Stack struct {
	list *list.List
	mut  sync.Mutex
}

// NewStack instantiates a new stack.
func NewStack() *Stack {
	return &Stack{
		list: list.New(),
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
