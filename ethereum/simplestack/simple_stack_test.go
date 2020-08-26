package simplestack

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	miniHeaderOne = &types.MiniHeader{
		Number:    big.NewInt(1),
		Hash:      common.HexToHash("0x1"),
		Parent:    common.HexToHash("0x0"),
		Timestamp: time.Now().UTC(),
	}
	miniHeaderTwo = &types.MiniHeader{
		Number:    big.NewInt(2),
		Hash:      common.HexToHash("0x2"),
		Parent:    common.HexToHash("0x1"),
		Timestamp: time.Now().UTC(),
	}
)

func TestSimpleStackPushPeekPop(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})
	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)

	expectedLen := 1
	miniHeaders := stack.PeekAll()
	assert.Len(t, miniHeaders, expectedLen)

	miniHeader := stack.Peek()
	assert.Equal(t, miniHeaders[0], miniHeader)

	expectedLen = 1
	miniHeaders = stack.PeekAll()
	assert.Len(t, miniHeaders, expectedLen)

	miniHeader = stack.Pop()
	assert.Equal(t, miniHeaders[0], miniHeader)

	expectedLen = 0
	miniHeaders = stack.PeekAll()
	assert.Len(t, miniHeaders, expectedLen)
}

func TestSimpleStackErrorIfPushTwoHeadersWithSameNumber(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})
	// Push miniHeaderOne
	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Push miniHeaderOne again
	err = stack.Push(miniHeaderOne)
	assert.Error(t, err)
}

func TestSimpleStackErrorIfResetWithoutCheckpointFirst(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})

	checkpointID := 123
	err := stack.Reset(checkpointID)
	require.Error(t, err)
}

func TestSimpleStackClear(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})

	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)

	miniHeader := stack.Peek()
	assert.Equal(t, miniHeaderOne, miniHeader)

	stack.Clear()

	miniHeader = stack.Peek()
	require.Nil(t, miniHeader)
}

func TestSimpleStackErrorIfResetWithOldCheckpoint(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})

	checkpointIDOne := stack.Checkpoint()
	checkpointIDTwo := stack.Checkpoint()

	err := stack.Reset(checkpointIDOne)
	require.Error(t, err)

	err = stack.Reset(checkpointIDTwo)
	require.NoError(t, err)
}

func TestSimpleStackCheckpoint(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})
	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)
	err = stack.Push(miniHeaderTwo)
	require.NoError(t, err)

	assert.Len(t, stack.updates, 2)

	stack.Checkpoint()

	assert.Len(t, stack.updates, 0)

	miniHeader := stack.Pop()
	assert.Equal(t, miniHeaderTwo, miniHeader)

	miniHeader = stack.Pop()
	assert.Equal(t, miniHeaderOne, miniHeader)

	assert.Len(t, stack.updates, 2)

	stack.Checkpoint()

	assert.Len(t, stack.updates, 0)
}

func TestSimpleStackCheckpointAfterSameHeaderPushedAndPopped(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})
	// Push miniHeaderOne
	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Pop miniHeaderOne
	miniHeader := stack.Pop()
	assert.Equal(t, miniHeaderOne, miniHeader)

	assert.Len(t, stack.miniHeaders, 0)
	assert.Len(t, stack.updates, 2)

	stack.Checkpoint()

	assert.Len(t, stack.updates, 0)
}

func TestSimpleStackCheckpointAfterSameHeaderPushedThenPoppedThenPushed(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})
	// Push miniHeaderOne
	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Pop miniHeaderOne
	miniHeader := stack.Pop()
	assert.Equal(t, miniHeaderOne, miniHeader)
	// Push miniHeaderOne again
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	assert.Len(t, stack.miniHeaders, 1)
	assert.Len(t, stack.updates, 3)

	stack.Checkpoint()

	assert.Len(t, stack.updates, 0)
}

func TestSimpleStackCheckpointThenReset(t *testing.T) {
	stack := New(10, []*types.MiniHeader{})

	checkpointID := stack.Checkpoint()

	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)

	assert.Len(t, stack.miniHeaders, 1)
	assert.Len(t, stack.updates, 1)

	err = stack.Reset(checkpointID)
	require.NoError(t, err)

	assert.Len(t, stack.miniHeaders, 0)
	assert.Len(t, stack.updates, 0)

	err = stack.Push(miniHeaderTwo)
	require.NoError(t, err)

	assert.Len(t, stack.miniHeaders, 1)
	assert.Len(t, stack.updates, 1)

	checkpointID = stack.Checkpoint()

	assert.Len(t, stack.miniHeaders, 1)
	assert.Len(t, stack.updates, 0)

	miniHeader := stack.Pop()
	assert.Equal(t, miniHeader, miniHeaderTwo)

	assert.Len(t, stack.miniHeaders, 0)
	assert.Len(t, stack.updates, 1)

	checkpointID = stack.Checkpoint()

	assert.Len(t, stack.miniHeaders, 0)
	assert.Len(t, stack.updates, 0)
}
