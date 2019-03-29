package blockwatch

import "container/list"

// Stack allows performing basic stack operations on a stack of MiniBlockHeaders.
type Stack struct {
	list *list.List
}

func NewStack() *Stack {
	return &Stack{
		list: list.New(),
	}
}

// Pop removes and returns the latest block on the block stack.
func (bs *Stack) Pop() *MiniBlockHeader {
	block := bs.list.Front()
	bs.list.Remove(block)
	return block.Value.(*MiniBlockHeader)
}

// Push pushes a block onto the block stack.
func (bs *Stack) Push(block *MiniBlockHeader) {
	bs.list.PushFront(block)
}

// Peek returns the latest block from the block stack without removing it.
func (bs *Stack) Peek() *MiniBlockHeader {
	block := bs.list.Front()
	if block == nil {
		return nil
	}
	return block.Value.(*MiniBlockHeader)
}

// Inspect returns all the blocks in the stack. This method should only be used for debugging
// and testing purposes. It is not performant.
func (bs *Stack) Inspect() []*MiniBlockHeader {
	blocks := []*MiniBlockHeader{}
	for e := bs.list.Back(); e != nil; e = e.Prev() {
		blocks = append(blocks, e.Value.(*MiniBlockHeader))
	}
	return blocks
}
