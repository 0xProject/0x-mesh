package core

import (
	"container/heap"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum/common"
)

// Ensure that ETHBackingHeap implements heap.Interface
var _ heap.Interface = &ETHBackingHeap{}

// An ETHBackingHeap is a min-heap of ETHBackings sorted by ETH per order.
type ETHBackingHeap []*meshdb.ETHBacking

func (h ETHBackingHeap) Len() int           { return len(h) }
func (h ETHBackingHeap) Less(i, j int) bool { return h[i].ETHPerOrder().Cmp(h[j].ETHPerOrder()) == -1 }
func (h ETHBackingHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ETHBackingHeap) Push(x interface{}) {
	*h = append(*h, x.(*meshdb.ETHBacking))
}

func (h *ETHBackingHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h ETHBackingHeap) Peek() *meshdb.ETHBacking {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}

// TODO(albrow): Could be optimized by using two data structures: a slice for
// sorting and a map for lookup.
func (h ETHBackingHeap) FindByMakerAddress(address common.Address) (*meshdb.ETHBacking, int) {
	for i, backing := range h {
		if backing.MakerAddress == address {
			return backing, i
		}
	}
	return nil, -1
}

func (h *ETHBackingHeap) UpdateByMakerAddress(address common.Address, diff int) {
	backing, index := h.FindByMakerAddress(address)
	backing.OrderCount += diff
	heap.Fix(h, index)
}

func (h *ETHBackingHeap) UpdateLowest(diff int) {
	lowest := heap.Pop(h).(*meshdb.ETHBacking)
	lowest.OrderCount += diff
	heap.Push(h, lowest)
}
