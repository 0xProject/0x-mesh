package blockwatch

// BlockStack allows performing basic stack operations on a stack of SuccinctBlocks.
type BlockStack struct {
	data []*SuccinctBlock
}

// Pop removes and returns the latest block on the block stack.
func (bs *BlockStack) Pop() *SuccinctBlock {
	i := len(bs.data) - 1
	block := bs.data[i]
	bs.data = bs.data[:i]
	return block
}

// Push pushes a block onto the block stack.
func (bs *BlockStack) Push(block *SuccinctBlock) {
	bs.data = append(bs.data, block)
}

// Peek returns the latest block from the block stack without removing it.
func (bs *BlockStack) Peek() *SuccinctBlock {
	if len(bs.data) == 0 {
		return nil
	}
	latestBlock := bs.data[len(bs.data)-1]
	return latestBlock
}
